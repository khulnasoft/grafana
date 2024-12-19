package secret

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	common "k8s.io/kube-openapi/pkg/common"

	secretV0Alpha1 "github.com/grafana/grafana/pkg/apis/secret/v0alpha1"
	"github.com/grafana/grafana/pkg/registry/apis/secret/reststorage"
	"github.com/grafana/grafana/pkg/services/apiserver/builder"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	secretstorage "github.com/grafana/grafana/pkg/storage/secret"
	"github.com/grafana/grafana/pkg/util"
)

var _ builder.APIGroupBuilder = (*SecretAPIBuilder)(nil)

type SecretAPIBuilder struct {
	secureValueStorage secretstorage.SecureValueStorage
	keeperStorage      secretstorage.KeeperStorage
}

func NewSecretAPIBuilder(secureValueStorage secretstorage.SecureValueStorage, keeperStorage secretstorage.KeeperStorage) *SecretAPIBuilder {
	return &SecretAPIBuilder{secureValueStorage, keeperStorage}
}

func RegisterAPIService(
	features featuremgmt.FeatureToggles,
	apiregistration builder.APIRegistrar,
	secureValueStorage secretstorage.SecureValueStorage,
	keeperStorage secretstorage.KeeperStorage,
) *SecretAPIBuilder {
	// Skip registration unless opting into experimental apis and the secrets management app platform flag.
	if !features.IsEnabledGlobally(featuremgmt.FlagGrafanaAPIServerWithExperimentalAPIs) ||
		!features.IsEnabledGlobally(featuremgmt.FlagSecretsManagementAppPlatform) {
		return nil
	}

	builder := NewSecretAPIBuilder(secureValueStorage, keeperStorage)
	apiregistration.RegisterAPI(builder)
	return builder
}

// GetGroupVersion returns the tuple of `group` and `version` for the API which uniquely identifies it.
func (b *SecretAPIBuilder) GetGroupVersion() schema.GroupVersion {
	return secretV0Alpha1.SchemeGroupVersion
}

// InstallSchema is called by the `apiserver` which exposes the defined kinds.
func (b *SecretAPIBuilder) InstallSchema(scheme *runtime.Scheme) error {
	secretV0Alpha1.AddKnownTypes(scheme, secretV0Alpha1.VERSION)

	// Link this version to the internal representation.
	// This is used for server-side-apply (PATCH), and avoids the error:
	// "no kind is registered for the type"
	secretV0Alpha1.AddKnownTypes(scheme, runtime.APIVersionInternal)

	// Internal Kubernetes metadata API. Presumably to display the available APIs?
	// e.g. http://localhost:3000/apis/secret.grafana.app/v0alpha1
	metav1.AddToGroupVersion(scheme, secretV0Alpha1.SchemeGroupVersion)

	// This sets the priority in case we have multiple versions.
	// By default Kubernetes will only let you use `kubectl get <resource>` with one version.
	// In case there are multiple versions, we'd need to pass the full path with the `--raw` flag.
	if err := scheme.SetVersionPriority(secretV0Alpha1.SchemeGroupVersion); err != nil {
		return fmt.Errorf("scheme set version priority: %w", err)
	}

	return nil
}

// UpdateAPIGroupInfo is called when creating a generic API server for this group of kinds.
func (b *SecretAPIBuilder) UpdateAPIGroupInfo(apiGroupInfo *genericapiserver.APIGroupInfo, opts builder.APIGroupOptions) error {
	secureValueResource := secretV0Alpha1.SecureValuesResourceInfo
	keeperResource := secretV0Alpha1.KeeperResourceInfo

	// rest.Storage is a generic interface for RESTful storage services.
	// The constructors need to at least implement this interface, but will most likely implement
	// other interfaces that equal to different operations like `get`, `list` and so on.
	secureRestStorage := map[string]rest.Storage{
		// Default path for `securevalue`.
		// The `reststorage.SecureValueRest` struct will implement interfaces for CRUDL operations on `securevalue`.
		secureValueResource.StoragePath(): reststorage.NewSecureValueRest(b.secureValueStorage, secureValueResource),

		// The `reststorage.KeeperRest` struct will implement interfaces for CRUDL operations on `keeper`.
		keeperResource.StoragePath(): reststorage.NewKeeperRest(b.keeperStorage, keeperResource),
	}

	apiGroupInfo.VersionedResourcesStorageMap[secretV0Alpha1.VERSION] = secureRestStorage
	return nil
}

// GetOpenAPIDefinitions, is this only for documentation?
func (b *SecretAPIBuilder) GetOpenAPIDefinitions() common.GetOpenAPIDefinitions {
	return secretV0Alpha1.GetOpenAPIDefinitions
}

// GetAuthorizer: [TODO] who can create secrets? must be multi-tenant first
func (b *SecretAPIBuilder) GetAuthorizer() authorizer.Authorizer {
	// This is TBD being defined with IAM. Test

	return nil // start with the default authorizer
}

// Register additional routes with the server.
func (b *SecretAPIBuilder) GetAPIRoutes() *builder.APIRoutes {
	return nil
}

// Validate is called in `Create`, `Update` and `Delete` REST funcs, if the body calls the argument `rest.ValidateObjectFunc`.
func (b *SecretAPIBuilder) Validate(ctx context.Context, a admission.Attributes, o admission.ObjectInterfaces) error {
	obj := a.GetObject()
	operation := a.GetOperation()

	if obj == nil || operation == admission.Connect {
		return nil // This is normal for sub-resource
	}

	groupKind := obj.GetObjectKind().GroupVersionKind().GroupKind()

	// Generic validations for all kinds. At this point the name+namespace must not be empty.
	if a.GetName() == "" {
		return apierrors.NewInvalid(
			groupKind,
			a.GetName(),
			field.ErrorList{field.Required(field.NewPath("metadata", "name"), "a `name` is required")},
		)
	}

	if a.GetNamespace() == "" {
		return apierrors.NewInvalid(
			groupKind,
			a.GetName(),
			field.ErrorList{field.Required(field.NewPath("metadata", "namespace"), "a `namespace` is required")},
		)
	}

	switch typedObj := obj.(type) {
	case *secretV0Alpha1.SecureValue:
		if errs := reststorage.ValidateSecureValue(typedObj, operation); len(errs) > 0 {
			return apierrors.NewInvalid(groupKind, a.GetName(), errs)
		}

		return nil
	case *secretV0Alpha1.Keeper:
		if errs := reststorage.ValidateKeeper(typedObj, operation); len(errs) > 0 {
			return apierrors.NewInvalid(groupKind, a.GetName(), errs)
		}

		return nil
	}

	return apierrors.NewBadRequest(fmt.Sprintf("unknown spec %T", obj))
}

func (b *SecretAPIBuilder) Mutate(ctx context.Context, a admission.Attributes, o admission.ObjectInterfaces) error {
	obj := a.GetObject()
	operation := a.GetOperation()

	if obj == nil || operation == admission.Connect {
		return nil // This is normal for sub-resource
	}

	// When creating a resource and the name is empty, we need to generate one.
	if operation == admission.Create && a.GetName() == "" {
		generatedName, err := util.GetRandomString(8)
		if err != nil {
			return fmt.Errorf("generate random string: %w", err)
		}

		switch typedObj := obj.(type) {
		case *secretV0Alpha1.SecureValue:
			optionalPrefix := typedObj.GenerateName
			if optionalPrefix == "" {
				optionalPrefix = "sv-"
			}

			typedObj.Name = optionalPrefix + generatedName

		case *secretV0Alpha1.Keeper:
			optionalPrefix := typedObj.GenerateName
			if optionalPrefix == "" {
				optionalPrefix = "kp-"
			}

			typedObj.Name = optionalPrefix + generatedName
		}
	}

	return nil
}
