package unified

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"gocloud.dev/blob/fileblob"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authnlib "github.com/grafana/authlib/authn"
	"github.com/grafana/authlib/types"

	infraDB "github.com/grafana/grafana/pkg/infra/db"
	"github.com/grafana/grafana/pkg/infra/tracing"
	"github.com/grafana/grafana/pkg/services/apiserver/options"
	"github.com/grafana/grafana/pkg/services/authn/grpcutils"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/storage/legacysql"
	"github.com/grafana/grafana/pkg/storage/unified/federated"
	"github.com/grafana/grafana/pkg/storage/unified/resource"
	"github.com/grafana/grafana/pkg/storage/unified/sql"
)

const resourceStoreAudience = "resourceStore"

var (
	// internal provider of the package level resource client
	pkgResourceClient resource.ResourceClient // Remove this
	ready             = make(chan struct{})
)

func GetResourceClient(ctx context.Context) resource.ResourceClient {
	<-ready
	return pkgResourceClient
}

type ClientService interface {
	GetResourceClient() (resource.ResourceClient, error)
}

func ProvideClientServiceImpl(cfg *setting.Cfg,
	features featuremgmt.FeatureToggles,
	db infraDB.DB,
	tracer tracing.Tracer,
	reg prometheus.Registerer,
	authzc types.AccessClient,
	docs resource.DocumentBuilderSupplier) *ClientServiceImpl {
	return &ClientServiceImpl{
		cfg:      cfg,
		features: features,
		db:       db,
		tracer:   tracer,
		reg:      reg,
		authzc:   authzc,
		docs:     docs,
	}
}

type ClientServiceImpl struct {
	cfg      *setting.Cfg
	features featuremgmt.FeatureToggles
	db       infraDB.DB
	tracer   tracing.Tracer
	reg      prometheus.Registerer
	authzc   types.AccessClient
	docs     resource.DocumentBuilderSupplier
	once     sync.Once
	client   resource.ResourceClient
	err      error
}

func (s *ClientServiceImpl) GetResourceClient() (resource.ResourceClient, error) {
	s.once.Do(func() {
		s.client, s.err = s.init()
	})
	return s.client, s.err
}

// This adds a UnifiedStorage client into the wire dependency tree
func (s *ClientServiceImpl) init() (resource.ResourceClient, error) {
	// See: apiserver.ApplyGrafanaConfig(cfg, features, o)
	apiserverCfg := s.cfg.SectionWithEnvOverrides("grafana-apiserver")
	client, err := newClient(options.StorageOptions{
		StorageType:  options.StorageType(apiserverCfg.Key("storage_type").MustString(string(options.StorageTypeUnified))),
		DataPath:     apiserverCfg.Key("storage_path").MustString(filepath.Join(s.cfg.DataPath, "grafana-apiserver")),
		Address:      apiserverCfg.Key("address").MustString(""), // client address
		BlobStoreURL: apiserverCfg.Key("blob_url").MustString(""),
	}, s.cfg, s.features, s.db, s.tracer, s.reg, s.authzc, s.docs)
	if err == nil {
		// Used to get the folder stats
		client = federated.NewFederatedClient(
			client, // The original
			legacysql.NewDatabaseProvider(s.db),
		)
	}

	// only set the package level restConfig once
	if pkgResourceClient == nil {
		pkgResourceClient = client
		close(ready)
	}
	return client, err
}

func newClient(opts options.StorageOptions,
	cfg *setting.Cfg,
	features featuremgmt.FeatureToggles,
	db infraDB.DB,
	tracer tracing.Tracer,
	reg prometheus.Registerer,
	authzc types.AccessClient,
	docs resource.DocumentBuilderSupplier,
) (resource.ResourceClient, error) {
	ctx := context.Background()
	switch opts.StorageType {
	case options.StorageTypeFile:
		if opts.DataPath == "" {
			opts.DataPath = filepath.Join(cfg.DataPath, "grafana-apiserver")
		}
		bucket, err := fileblob.OpenBucket(filepath.Join(opts.DataPath, "resource"), &fileblob.Options{
			CreateDir: true,
			Metadata:  fileblob.MetadataDontWrite, // skip
		})
		if err != nil {
			return nil, err
		}
		backend, err := resource.NewCDKBackend(ctx, resource.CDKBackendOptions{
			Bucket: bucket,
		})
		if err != nil {
			return nil, err
		}
		server, err := resource.NewResourceServer(resource.ResourceServerOptions{
			Backend: backend,
			Blob: resource.BlobConfig{
				URL: opts.BlobStoreURL,
			},
		})
		if err != nil {
			return nil, err
		}
		return resource.NewLocalResourceClient(server), nil

	case options.StorageTypeUnifiedGrpc:
		if opts.Address == "" {
			return nil, fmt.Errorf("expecting address for storage_type: %s", opts.StorageType)
		}

		// Create a connection to the gRPC server
		conn, err := grpc.NewClient(opts.Address,
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, err
		}

		// Create a client instance
		client, err := newResourceClient(conn, cfg, features, tracer)
		if err != nil {
			return nil, err
		}
		return client, nil

	// Use the local SQL
	default:
		server, err := sql.NewResourceServer(ctx, db, cfg, features, docs, tracer, reg, authzc)
		if err != nil {
			return nil, err
		}
		return resource.NewLocalResourceClient(server), nil
	}
}

func clientCfgMapping(clientCfg *grpcutils.GrpcClientConfig) authnlib.GrpcClientConfig {
	return authnlib.GrpcClientConfig{
		TokenClientConfig: &authnlib.TokenExchangeConfig{
			Token:            clientCfg.Token,
			TokenExchangeURL: clientCfg.TokenExchangeURL,
		},
		TokenRequest: &authnlib.TokenExchangeRequest{
			Namespace: clientCfg.TokenNamespace,
			Audiences: []string{resourceStoreAudience},
		},
	}
}

func newResourceClient(conn *grpc.ClientConn, cfg *setting.Cfg, features featuremgmt.FeatureToggles, tracer tracing.Tracer) (resource.ResourceClient, error) {
	if !features.IsEnabledGlobally(featuremgmt.FlagAppPlatformGrpcClientAuth) {
		return resource.NewLegacyResourceClient(conn), nil
	}
	if cfg.StackID == "" {
		return resource.NewGRPCResourceClient(tracer, conn)
	}

	grpcClientCfg := grpcutils.ReadGrpcClientConfig(cfg)

	return resource.NewCloudResourceClient(tracer, conn, clientCfgMapping(grpcClientCfg), cfg.Env == setting.Dev)
}
