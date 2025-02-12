package metadata

import (
	"context"
	"fmt"

	claims "github.com/grafana/authlib/types"
	secretv0alpha1 "github.com/grafana/grafana/pkg/apis/secret/v0alpha1"
	"github.com/grafana/grafana/pkg/infra/db"
	"github.com/grafana/grafana/pkg/registry/apis/secret/contracts"
	"github.com/grafana/grafana/pkg/registry/apis/secret/secretkeeper"
	keepertypes "github.com/grafana/grafana/pkg/registry/apis/secret/secretkeeper/types"
	"github.com/grafana/grafana/pkg/registry/apis/secret/xkube"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/services/sqlstore"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/storage/secret/migrator"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	"k8s.io/apimachinery/pkg/labels"
)

func ProvideSecureValueStorage(db db.DB, cfg *setting.Cfg, features featuremgmt.FeatureToggles, keeperService secretkeeper.Service) (contracts.SecureValueStorage, error) {
	if !features.IsEnabledGlobally(featuremgmt.FlagGrafanaAPIServerWithExperimentalAPIs) ||
		!features.IsEnabledGlobally(featuremgmt.FlagSecretsManagementAppPlatform) {
		return &secureValueStorage{}, nil
	}

	if err := migrator.MigrateSecretSQL(db.GetEngine(), cfg); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	keepers, err := keeperService.GetKeepers()
	if err != nil {
		return nil, fmt.Errorf("failed to get keepers: %w", err)
	}

	return &secureValueStorage{db: db, keepers: keepers}, nil
}

// secureValueStorage is the actual implementation of the secure value (metadata) storage.
type secureValueStorage struct {
	db      db.DB
	keepers map[keepertypes.KeeperType]keepertypes.Keeper
}

func (s *secureValueStorage) Create(ctx context.Context, sv *secretv0alpha1.SecureValue) (*secretv0alpha1.SecureValue, error) {
	authInfo, ok := claims.AuthInfoFrom(ctx)
	if !ok {
		return nil, fmt.Errorf("missing auth info in context")
	}

	// Store in keeper.
	externalID, err := s.storeInKeeper(ctx, sv)
	if err != nil {
		return nil, fmt.Errorf("failed to store in keeper: %w", err)
	}

	// From this point on, we should not have a need to read value.
	sv.Spec.Value = ""

	row, err := toCreateRow(sv, authInfo.GetUID(), externalID.String())
	if err != nil {
		return nil, fmt.Errorf("to create row: %w", err)
	}

	err = s.db.WithDbSession(ctx, func(sess *sqlstore.DBSession) error {
		if _, err := sess.Insert(row); err != nil {
			return fmt.Errorf("insert row: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("db failure: %w", err)
	}

	createdSecureValue, err := row.toKubernetes()
	if err != nil {
		return nil, fmt.Errorf("convert to kubernetes object: %w", err)
	}

	return createdSecureValue, nil
}

func (s *secureValueStorage) Read(ctx context.Context, namespace xkube.Namespace, name string) (*secretv0alpha1.SecureValue, error) {
	_, ok := claims.AuthInfoFrom(ctx)
	if !ok {
		return nil, fmt.Errorf("missing auth info in context")
	}

	row := &secureValueDB{Name: name, Namespace: namespace.String()}

	err := s.db.WithDbSession(ctx, func(sess *sqlstore.DBSession) error {
		found, err := sess.Get(row)
		if err != nil {
			return fmt.Errorf("could not get row: %w", err)
		}

		if !found {
			return contracts.ErrSecureValueNotFound
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("db failure: %w", err)
	}

	secureValue, err := row.toKubernetes()
	if err != nil {
		return nil, fmt.Errorf("convert to kubernetes object: %w", err)
	}

	return secureValue, nil
}

func (s *secureValueStorage) Update(ctx context.Context, newSecureValue *secretv0alpha1.SecureValue) (*secretv0alpha1.SecureValue, error) {
	authInfo, ok := claims.AuthInfoFrom(ctx)
	if !ok {
		return nil, fmt.Errorf("missing auth info in context")
	}

	currentRow := &secureValueDB{Name: newSecureValue.Name, Namespace: newSecureValue.Namespace}

	err := s.db.WithDbSession(ctx, func(sess *sqlstore.DBSession) error {
		found, err := sess.Get(currentRow)
		if err != nil {
			return fmt.Errorf("could not get row: %w", err)
		}

		if !found {
			return contracts.ErrSecureValueNotFound
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("db failure: %w", err)
	}

	// Update in keeper at value updates.
	err = s.updateInKeeper(ctx, currentRow, newSecureValue)
	if err != nil {
		return nil, fmt.Errorf("failed to update in keeper: %w", err)
	}

	// From this point on, we should not have a need to read value.
	newSecureValue.Spec.Value = ""

	newRow, err := toUpdateRow(currentRow, newSecureValue, authInfo.GetUID(), currentRow.ExternalID)
	if err != nil {
		return nil, fmt.Errorf("to update row: %w", err)
	}

	err = s.db.WithDbSession(ctx, func(sess *sqlstore.DBSession) error {
		cond := &secureValueDB{Name: newSecureValue.Name, Namespace: newSecureValue.Namespace}

		if _, err := sess.Update(newRow, cond); err != nil {
			return fmt.Errorf("update row: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("db failure: %w", err)
	}

	secureValue, err := newRow.toKubernetes()
	if err != nil {
		return nil, fmt.Errorf("convert to kubernetes object: %w", err)
	}

	return secureValue, nil
}

func (s *secureValueStorage) Delete(ctx context.Context, namespace xkube.Namespace, name string) error {
	_, ok := claims.AuthInfoFrom(ctx)
	if !ok {
		return fmt.Errorf("missing auth info in context")
	}

	// Delete from the keeper.
	// TODO: log any errors, but do not inform the caller what the errors are.
	_ = s.deleteFromKeeper(ctx, namespace, name)

	// TODO: do we need to delete by GUID? name+namespace is a unique index. It would avoid doing a fetch.
	row := &secureValueDB{Name: name, Namespace: namespace.String()}

	err := s.db.WithDbSession(ctx, func(sess *sqlstore.DBSession) error {
		// TODO: because this is a securevalue, do we care to inform the caller if a row was delete (existed) or not?
		if _, err := sess.Delete(row); err != nil {
			return fmt.Errorf("delete row: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("db failure: %w", err)
	}

	return nil
}

func (s *secureValueStorage) List(ctx context.Context, namespace xkube.Namespace, options *internalversion.ListOptions) (*secretv0alpha1.SecureValueList, error) {
	_, ok := claims.AuthInfoFrom(ctx)
	if !ok {
		return nil, fmt.Errorf("missing auth info in context")
	}

	labelSelector := options.LabelSelector
	if labelSelector == nil {
		labelSelector = labels.Everything()
	}

	secureValueRows := make([]*secureValueDB, 0)

	err := s.db.WithDbSession(ctx, func(sess *sqlstore.DBSession) error {
		cond := &secureValueDB{Namespace: namespace.String()}

		if err := sess.Find(&secureValueRows, cond); err != nil {
			return fmt.Errorf("find rows: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("db failure: %w", err)
	}

	secureValues := make([]secretv0alpha1.SecureValue, 0, len(secureValueRows))

	for _, row := range secureValueRows {
		secureValue, err := row.toKubernetes()
		if err != nil {
			return nil, fmt.Errorf("convert to kubernetes object: %w", err)
		}

		if labelSelector.Matches(labels.Set(secureValue.Labels)) {
			secureValues = append(secureValues, *secureValue)
		}
	}

	return &secretv0alpha1.SecureValueList{
		Items: secureValues,
	}, nil
}

func (s *secureValueStorage) readInternal(ctx context.Context, namespace xkube.Namespace, name string) (*secureValueDB, error) {
	_, ok := claims.AuthInfoFrom(ctx)
	if !ok {
		return nil, fmt.Errorf("missing auth info in context")
	}

	row := &secureValueDB{Namespace: namespace.String(), Name: name}

	err := s.db.WithDbSession(ctx, func(sess *sqlstore.DBSession) error {
		found, err := sess.Get(row)
		if err != nil {
			return fmt.Errorf("could not get row: %w", err)
		}

		if !found {
			return contracts.ErrSecureValueNotFound
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("db failure: %w", err)
	}

	return row, nil
}

func (s *secureValueStorage) storeInKeeper(ctx context.Context, sv *secretv0alpha1.SecureValue) (keepertypes.ExternalID, error) {
	// TODO: Implement store by ref
	if sv.Spec.Ref != "" {
		return "", fmt.Errorf("store by ref in keeper")
	}

	// Check if keeper is default.
	if sv.Spec.Keeper == keepertypes.DefaultKeeper {
		keeper, exists := s.keepers[keepertypes.SQLKeeperType]
		if !exists {
			return "", fmt.Errorf("could not find default keeper")
		}
		externalID, err := keeper.Store(ctx, nil, sv.Namespace, string(sv.Spec.Value))
		if err != nil {
			return "", fmt.Errorf("failed to store in default keeper: %w", err)
		}
		return externalID, err
	}

	// Load keeper config from metadata store.
	keeperType, keeperConfig, err := s.getKeeperConfig(ctx, sv.Spec.Keeper, sv.Namespace)
	if err != nil {
		return "", fmt.Errorf("get keeper config: %w", err)
	}

	// Store in keeper.
	keeper, ok := s.keepers[keeperType]
	if !ok {
		return "", fmt.Errorf("could not find keeper: %s", keeperType)
	}

	return keeper.Store(ctx, keeperConfig, sv.Namespace, string(sv.Spec.Value))
}

func (s *secureValueStorage) updateInKeeper(ctx context.Context, currRow *secureValueDB, newSV *secretv0alpha1.SecureValue) error {
	// TODO: Implement store by ref
	if newSV.Spec.Ref != "" {
		return fmt.Errorf("store by ref in keeper")
	}

	// Check if an update is needed.
	if newSV.Spec.Value == "" {
		return nil
	}

	if currRow.Keeper != newSV.Spec.Keeper {
		return fmt.Errorf("keeper change not supported")
	}

	// Check if keeper is default.
	if currRow.Keeper == keepertypes.DefaultKeeper {
		keeper, exists := s.keepers[keepertypes.SQLKeeperType]
		if !exists {
			return fmt.Errorf("could not find default keeper")
		}
		return keeper.Update(ctx, nil, currRow.Namespace, keepertypes.ExternalID(currRow.ExternalID), string(newSV.Spec.Value))
	}

	// Load keeper config from metadata store.
	keeperType, keeperConfig, err := s.getKeeperConfig(ctx, currRow.Keeper, currRow.Namespace)
	if err != nil {
		return fmt.Errorf("get keeper config: %w", err)
	}

	// Store in keeper.
	keeper, ok := s.keepers[keeperType]
	if !ok {
		return fmt.Errorf("could not find keeper: %s", keeperType)
	}

	return keeper.Update(ctx, keeperConfig, currRow.Namespace, keepertypes.ExternalID(currRow.ExternalID), string(newSV.Spec.Value))
}

func (s *secureValueStorage) deleteFromKeeper(ctx context.Context, namespace xkube.Namespace, name string) error {
	sv, err := s.readInternal(ctx, namespace, name)
	if err != nil {
		return fmt.Errorf("read securevalue: %w", err)
	}

	// Check if keeper is default.
	if sv.Keeper == keepertypes.DefaultKeeper {
		keeper, exists := s.keepers[keepertypes.SQLKeeperType]
		if !exists {
			return fmt.Errorf("could not find default keeper")
		}
		return keeper.Delete(ctx, nil, namespace.String(), keepertypes.ExternalID(sv.ExternalID))
	}

	// Load keeper config from metadata store.
	keeperType, keeperConfig, err := s.getKeeperConfig(ctx, sv.Keeper, namespace.String())
	if err != nil {
		return fmt.Errorf("get keeper config: %w", err)
	}

	// Delete from keeper.
	keeper, ok := s.keepers[keeperType]
	if !ok {
		return fmt.Errorf("could not find keeper: %s", keeperType)
	}
	return keeper.Delete(ctx, keeperConfig, namespace.String(), keepertypes.ExternalID(sv.ExternalID))
}

// this feels like it should belong to keeper store
func (s *secureValueStorage) getKeeperConfig(ctx context.Context, keeperName string, namespace string) (keepertypes.KeeperType, secretv0alpha1.KeeperConfig, error) {
	kp := &keeperDB{Name: keeperName, Namespace: namespace}
	err := s.db.WithDbSession(ctx, func(sess *sqlstore.DBSession) error {
		found, err := sess.Get(kp)
		if err != nil {
			return fmt.Errorf("failed to get row: %w", err)
		}
		if !found {
			return contracts.ErrKeeperNotFound
		}

		return nil
	})
	if err != nil {
		return "", nil, fmt.Errorf("db failure: %w", err)
	}

	keeperConfig := toProvider(kp.Type, kp.Payload)
	keeperType := keepertypes.SQLKeeperType // TODO: mapping between keeperDB.Type and KeeperType

	return keeperType, keeperConfig, nil
}
