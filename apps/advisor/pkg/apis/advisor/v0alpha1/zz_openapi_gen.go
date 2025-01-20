//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by grafana-app-sdk. DO NOT EDIT.

package v0alpha1

import (
	common "k8s.io/kube-openapi/pkg/common"
	spec "k8s.io/kube-openapi/pkg/validation/spec"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/grafana/grafana/apps/advisor/pkg/apis/advisor/v0alpha1.Check":                           schema_pkg_apis_advisor_v0alpha1_Check(ref),
		"github.com/grafana/grafana/apps/advisor/pkg/apis/advisor/v0alpha1.CheckList":                       schema_pkg_apis_advisor_v0alpha1_CheckList(ref),
		"github.com/grafana/grafana/apps/advisor/pkg/apis/advisor/v0alpha1.CheckSpec":                       schema_pkg_apis_advisor_v0alpha1_CheckSpec(ref),
		"github.com/grafana/grafana/apps/advisor/pkg/apis/advisor/v0alpha1.CheckStatus":                     schema_pkg_apis_advisor_v0alpha1_CheckStatus(ref),
		"github.com/grafana/grafana/apps/advisor/pkg/apis/advisor/v0alpha1.CheckV0alpha1StatusReport":       schema_pkg_apis_advisor_v0alpha1_CheckV0alpha1StatusReport(ref),
		"github.com/grafana/grafana/apps/advisor/pkg/apis/advisor/v0alpha1.CheckV0alpha1StatusReportErrors": schema_pkg_apis_advisor_v0alpha1_CheckV0alpha1StatusReportErrors(ref),
		"github.com/grafana/grafana/apps/advisor/pkg/apis/advisor/v0alpha1.CheckstatusOperatorState":        schema_pkg_apis_advisor_v0alpha1_CheckstatusOperatorState(ref),
	}
}

func schema_pkg_apis_advisor_v0alpha1_Check(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("github.com/grafana/grafana/apps/advisor/pkg/apis/advisor/v0alpha1.CheckSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("github.com/grafana/grafana/apps/advisor/pkg/apis/advisor/v0alpha1.CheckStatus"),
						},
					},
				},
				Required: []string{"metadata", "spec", "status"},
			},
		},
		Dependencies: []string{
			"github.com/grafana/grafana/apps/advisor/pkg/apis/advisor/v0alpha1.CheckSpec", "github.com/grafana/grafana/apps/advisor/pkg/apis/advisor/v0alpha1.CheckStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_advisor_v0alpha1_CheckList(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("k8s.io/apimachinery/pkg/apis/meta/v1.ListMeta"),
						},
					},
					"items": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: map[string]interface{}{},
										Ref:     ref("github.com/grafana/grafana/apps/advisor/pkg/apis/advisor/v0alpha1.Check"),
									},
								},
							},
						},
					},
				},
				Required: []string{"metadata", "items"},
			},
		},
		Dependencies: []string{
			"github.com/grafana/grafana/apps/advisor/pkg/apis/advisor/v0alpha1.Check", "k8s.io/apimachinery/pkg/apis/meta/v1.ListMeta"},
	}
}

func schema_pkg_apis_advisor_v0alpha1_CheckSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"data": {
						SchemaProps: spec.SchemaProps{
							Description: "Generic data input that a check can receive",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: "",
										Type:    []string{"string"},
										Format:  "",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func schema_pkg_apis_advisor_v0alpha1_CheckStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"report": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("github.com/grafana/grafana/apps/advisor/pkg/apis/advisor/v0alpha1.CheckV0alpha1StatusReport"),
						},
					},
					"operatorStates": {
						SchemaProps: spec.SchemaProps{
							Description: "operatorStates is a map of operator ID to operator state evaluations. Any operator which consumes this kind SHOULD add its state evaluation information to this field.",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: map[string]interface{}{},
										Ref:     ref("github.com/grafana/grafana/apps/advisor/pkg/apis/advisor/v0alpha1.CheckstatusOperatorState"),
									},
								},
							},
						},
					},
					"additionalFields": {
						SchemaProps: spec.SchemaProps{
							Description: "additionalFields is reserved for future use",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"object"},
										Format: "",
									},
								},
							},
						},
					},
				},
				Required: []string{"report"},
			},
		},
		Dependencies: []string{
			"github.com/grafana/grafana/apps/advisor/pkg/apis/advisor/v0alpha1.CheckV0alpha1StatusReport", "github.com/grafana/grafana/apps/advisor/pkg/apis/advisor/v0alpha1.CheckstatusOperatorState"},
	}
}

func schema_pkg_apis_advisor_v0alpha1_CheckV0alpha1StatusReport(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"count": {
						SchemaProps: spec.SchemaProps{
							Description: "Number of elements analyzed",
							Default:     0,
							Type:        []string{"integer"},
							Format:      "int64",
						},
					},
					"errors": {
						SchemaProps: spec.SchemaProps{
							Description: "List of errors",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: map[string]interface{}{},
										Ref:     ref("github.com/grafana/grafana/apps/advisor/pkg/apis/advisor/v0alpha1.CheckV0alpha1StatusReportErrors"),
									},
								},
							},
						},
					},
				},
				Required: []string{"count", "errors"},
			},
		},
		Dependencies: []string{
			"github.com/grafana/grafana/apps/advisor/pkg/apis/advisor/v0alpha1.CheckV0alpha1StatusReportErrors"},
	}
}

func schema_pkg_apis_advisor_v0alpha1_CheckV0alpha1StatusReportErrors(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"type": {
						SchemaProps: spec.SchemaProps{
							Description: "Investigation or Action recommended",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"reason": {
						SchemaProps: spec.SchemaProps{
							Description: "Reason for the error",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"action": {
						SchemaProps: spec.SchemaProps{
							Description: "Action to take",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"type", "reason", "action"},
			},
		},
	}
}

func schema_pkg_apis_advisor_v0alpha1_CheckstatusOperatorState(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"lastEvaluation": {
						SchemaProps: spec.SchemaProps{
							Description: "lastEvaluation is the ResourceVersion last evaluated",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"state": {
						SchemaProps: spec.SchemaProps{
							Description: "state describes the state of the lastEvaluation. It is limited to three possible states for machine evaluation.",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"descriptiveState": {
						SchemaProps: spec.SchemaProps{
							Description: "descriptiveState is an optional more descriptive state field which has no requirements on format",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"details": {
						SchemaProps: spec.SchemaProps{
							Description: "details contains any extra information that is operator-specific",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"object"},
										Format: "",
									},
								},
							},
						},
					},
				},
				Required: []string{"lastEvaluation", "state"},
			},
		},
	}
}
