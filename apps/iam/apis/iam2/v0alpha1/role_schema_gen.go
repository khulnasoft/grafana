//
// Code generated by grafana-app-sdk. DO NOT EDIT.
//

package v0alpha1

import (
	"github.com/grafana/grafana-app-sdk/resource"
)

// schema is unexported to prevent accidental overwrites
var (
	schemaRole = resource.NewSimpleSchema("iam2.grafana.app", "v0alpha1", &Role{}, &RoleList{}, resource.WithKind("Role"),
		resource.WithPlural("roles"), resource.WithScope(resource.NamespacedScope))
	kindRole = resource.Kind{
		Schema: schemaRole,
		Codecs: map[resource.KindEncoding]resource.Codec{
			resource.KindEncodingJSON: &RoleJSONCodec{},
		},
	}
)

// Kind returns a resource.Kind for this Schema with a JSON codec
func RoleKind() resource.Kind {
	return kindRole
}

// Schema returns a resource.SimpleSchema representation of Role
func RoleSchema() *resource.SimpleSchema {
	return schemaRole
}

// Interface compliance checks
var _ resource.Schema = kindRole
