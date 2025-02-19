package unified

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"gocloud.dev/blob/fileblob"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authnlib "github.com/grafana/authlib/authn"
	"github.com/grafana/authlib/types"
	"github.com/grafana/dskit/grpcclient"
	"github.com/grafana/dskit/middleware"

	infraDB "github.com/grafana/grafana/pkg/infra/db"
	"github.com/grafana/grafana/pkg/infra/tracing"
	"github.com/grafana/grafana/pkg/services/apiserver/options"
	"github.com/grafana/grafana/pkg/services/authn/grpcutils"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/storage/legacysql"
	"github.com/grafana/grafana/pkg/storage/unified/federated"
	"github.com/grafana/grafana/pkg/storage/unified/resource"
	"github.com/grafana/grafana/pkg/storage/unified/search"
	"github.com/grafana/grafana/pkg/storage/unified/sql"
)

const resourceStoreAudience = "resourceStore"

type Options struct {
	Cfg      *setting.Cfg
	Features featuremgmt.FeatureToggles
	DB       infraDB.DB
	Tracer   tracing.Tracer
	Reg      prometheus.Registerer
	Authzc   types.AccessClient
	Docs     resource.DocumentBuilderSupplier
}

type ClientMetrics struct {
	requestDuration *prometheus.HistogramVec
}

// This adds a UnifiedStorage client into the wire dependency tree
func ProvideUnifiedStorageClient(opts *Options) (resource.ResourceClient, error) {
	// See: apiserver.ApplyGrafanaConfig(cfg, features, o)
	apiserverCfg := opts.Cfg.SectionWithEnvOverrides("grafana-apiserver")
	client, err := newClient(options.StorageOptions{
		StorageType:  options.StorageType(apiserverCfg.Key("storage_type").MustString(string(options.StorageTypeUnified))),
		DataPath:     apiserverCfg.Key("storage_path").MustString(filepath.Join(opts.Cfg.DataPath, "grafana-apiserver")),
		Address:      apiserverCfg.Key("address").MustString(""), // client address
		BlobStoreURL: apiserverCfg.Key("blob_url").MustString(""),
	}, opts.Cfg, opts.Features, opts.DB, opts.Tracer, opts.Reg, opts.Authzc, opts.Docs)
	if err == nil {
		// Used to get the folder stats
		client = federated.NewFederatedClient(
			client, // The original
			legacysql.NewDatabaseProvider(opts.DB),
		)
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

		// Create a connection to the gRPC server.
		conn, err := grpcConn(opts.Address, reg)
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
		searchOptions, err := search.NewSearchOptions(features, cfg, tracer, docs, reg)
		if err != nil {
			return nil, err
		}
		server, err := sql.NewResourceServer(db, cfg, tracer, reg, authzc, searchOptions)
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
	return resource.NewRemoteResourceClient(tracer, conn, clientCfgMapping(grpcutils.ReadGrpcClientConfig(cfg)), cfg.Env == setting.Dev)
}

// grpcConn creates a new gRPC connection to the provided address.
func grpcConn(address string, reg prometheus.Registerer) (*grpc.ClientConn, error) {
	// This works for now as the Provide function is only called once during startup.
	// We might eventually want to tie this factory to a struct for more runtime control.
	metrics := ClientMetrics{
		requestDuration: promauto.With(reg).NewHistogramVec(prometheus.HistogramOpts{
			Name:    "resource_server_client_request_duration_seconds",
			Help:    "Time spent executing requests to the resource server.",
			Buckets: prometheus.ExponentialBuckets(0.008, 4, 7),
		}, []string{"operation", "status_code"}),
	}

	// Report gRPC status code errors as labels.
	var instrumentationOptions []middleware.InstrumentationOption
	instrumentationOptions = append(instrumentationOptions, middleware.ReportGRPCStatusOption)
	unary, stream := grpcclient.Instrument(metrics.requestDuration, instrumentationOptions...)

	// We can later pass in the gRPC config here, i.e. to set MaxRecvMsgSize etc.
	cfg := grpcclient.Config{}
	opts, err := cfg.DialOption(unary, stream)
	if err != nil {
		return nil, fmt.Errorf("could not instrument grpc client: %w", err)
	}

	opts = append(opts, grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// Use round_robin to balance requests more evenly over the available Storage server.
	opts = append(opts, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))

	// Disable looking up service config from TXT DNS records.
	// This reduces the number of requests made to the DNS servers.
	opts = append(opts, grpc.WithDisableServiceConfig())

	// Create a connection to the gRPC server
	return grpc.NewClient(address, opts...)
}
