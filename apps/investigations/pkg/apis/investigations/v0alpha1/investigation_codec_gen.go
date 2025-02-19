//
// Code generated by grafana-app-sdk. DO NOT EDIT.
//

package v0alpha1

import (
	"encoding/json"
	"io"

	"github.com/grafana/grafana-app-sdk/resource"
)

// InvestigationJSONCodec is an implementation of resource.Codec for kubernetes JSON encoding
type InvestigationJSONCodec struct{}

// Read reads JSON-encoded bytes from `reader` and unmarshals them into `into`
func (*InvestigationJSONCodec) Read(reader io.Reader, into resource.Object) error {
	return json.NewDecoder(reader).Decode(into)
}

// Write writes JSON-encoded bytes into `writer` marshaled from `from`
func (*InvestigationJSONCodec) Write(writer io.Writer, from resource.Object) error {
	return json.NewEncoder(writer).Encode(from)
}

// Interface compliance checks
var _ resource.Codec = &InvestigationJSONCodec{}
