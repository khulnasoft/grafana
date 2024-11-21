//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// SPDX-License-Identifier: AGPL-3.0-only

// Code generated by openapi-gen. DO NOT EDIT.

package v0alpha1

import (
	common "k8s.io/kube-openapi/pkg/common"
	spec "k8s.io/kube-openapi/pkg/validation/spec"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1.EditingOptions":         schema_pkg_apis_provisioning_v0alpha1_EditingOptions(ref),
		"github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1.GitHubRepositoryConfig": schema_pkg_apis_provisioning_v0alpha1_GitHubRepositoryConfig(ref),
		"github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1.HelloWorld":             schema_pkg_apis_provisioning_v0alpha1_HelloWorld(ref),
		"github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1.LocalRepositoryConfig":  schema_pkg_apis_provisioning_v0alpha1_LocalRepositoryConfig(ref),
		"github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1.Repository":             schema_pkg_apis_provisioning_v0alpha1_Repository(ref),
		"github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1.RepositoryList":         schema_pkg_apis_provisioning_v0alpha1_RepositoryList(ref),
		"github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1.RepositorySpec":         schema_pkg_apis_provisioning_v0alpha1_RepositorySpec(ref),
		"github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1.RepositoryStatus":       schema_pkg_apis_provisioning_v0alpha1_RepositoryStatus(ref),
		"github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1.ResourceWrapper":        schema_pkg_apis_provisioning_v0alpha1_ResourceWrapper(ref),
		"github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1.S3RepositoryConfig":     schema_pkg_apis_provisioning_v0alpha1_S3RepositoryConfig(ref),
		"github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1.WebhookResponse":        schema_pkg_apis_provisioning_v0alpha1_WebhookResponse(ref),
	}
}

func schema_pkg_apis_provisioning_v0alpha1_EditingOptions(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"create": {
						SchemaProps: spec.SchemaProps{
							Description: "End users can create new files in the remote file system",
							Default:     false,
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"update": {
						SchemaProps: spec.SchemaProps{
							Description: "End users can update existing files in the remote file system",
							Default:     false,
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"delete": {
						SchemaProps: spec.SchemaProps{
							Description: "End users can delete existing files in the remote file system",
							Default:     false,
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
				},
				Required: []string{"create", "update", "delete"},
			},
		},
	}
}

func schema_pkg_apis_provisioning_v0alpha1_GitHubRepositoryConfig(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"owner": {
						SchemaProps: spec.SchemaProps{
							Description: "The owner of the repository (e.g. example in `example/test` or `https://github.com/example/test`).",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"repository": {
						SchemaProps: spec.SchemaProps{
							Description: "The name of the repository (e.g. test in `example/test` or `https://github.com/example/test`).",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"branch": {
						SchemaProps: spec.SchemaProps{
							Description: "The branch to use in the repository. By default, this is the main branch.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"token": {
						SchemaProps: spec.SchemaProps{
							Description: "Token for accessing the repository.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"webhookURL": {
						SchemaProps: spec.SchemaProps{
							Description: "WebhookURL is the URL to send webhooks events to. By default, the system will generate a URL for you but you can use this one to run grafana locally and test the webhooks.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"webhookSecret": {
						SchemaProps: spec.SchemaProps{
							Description: "WebhookSecret is the secret used to validate incoming webhooks.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"branchWorkflow": {
						SchemaProps: spec.SchemaProps{
							Description: "Whether we should commit to change branches and use a Pull Request flow to achieve this. By default, this is false (i.e. we will commit straight to the main branch).",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"generateDashboardPreviews": {
						SchemaProps: spec.SchemaProps{
							Description: "Whether we should show dashboard previews in the pull requests caused by the BranchWorkflow option. By default, this is false (i.e. we will not create previews). This option is a no-op if BranchWorkflow is `false` or default.",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
				},
			},
		},
	}
}

func schema_pkg_apis_provisioning_v0alpha1_HelloWorld(ref common.ReferenceCallback) common.OpenAPIDefinition {
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
					"whom": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
				},
			},
		},
	}
}

func schema_pkg_apis_provisioning_v0alpha1_LocalRepositoryConfig(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"path": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
				},
			},
		},
	}
}

func schema_pkg_apis_provisioning_v0alpha1_Repository(ref common.ReferenceCallback) common.OpenAPIDefinition {
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
							Ref:     ref("github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1.RepositorySpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1.RepositoryStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1.RepositorySpec", "github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1.RepositoryStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_provisioning_v0alpha1_RepositoryList(ref common.ReferenceCallback) common.OpenAPIDefinition {
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
										Ref:     ref("github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1.Repository"),
									},
								},
							},
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1.Repository", "k8s.io/apimachinery/pkg/apis/meta/v1.ListMeta"},
	}
}

func schema_pkg_apis_provisioning_v0alpha1_RepositorySpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"title": {
						SchemaProps: spec.SchemaProps{
							Description: "Describe the feature toggle",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"description": {
						SchemaProps: spec.SchemaProps{
							Description: "Describe the feature toggle",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"folder": {
						SchemaProps: spec.SchemaProps{
							Description: "The folder that is backed by the repository. The value is a reference to the Kubernetes metadata name of the folder in the same namespace.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"editing": {
						SchemaProps: spec.SchemaProps{
							Description: "Edit options within the repository",
							Default:     map[string]interface{}{},
							Ref:         ref("github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1.EditingOptions"),
						},
					},
					"type": {
						SchemaProps: spec.SchemaProps{
							Description: "The repository type.  When selected oneOf the values below should be non-nil\n\nPossible enum values:\n - `\"github\"`\n - `\"local\"`\n - `\"s3\"`",
							Default:     "",
							Type:        []string{"string"},
							Format:      "",
							Enum:        []interface{}{"github", "local", "s3"},
						},
					},
					"local": {
						SchemaProps: spec.SchemaProps{
							Description: "The repository on the local file system. Mutually exclusive with s3 and github.",
							Ref:         ref("github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1.LocalRepositoryConfig"),
						},
					},
					"s3": {
						SchemaProps: spec.SchemaProps{
							Description: "The repository in an S3 bucket. Mutually exclusive with local and github.",
							Ref:         ref("github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1.S3RepositoryConfig"),
						},
					},
					"github": {
						SchemaProps: spec.SchemaProps{
							Description: "The repository on GitHub. Mutually exclusive with local and s3.",
							Ref:         ref("github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1.GitHubRepositoryConfig"),
						},
					},
				},
				Required: []string{"title", "editing", "type"},
			},
		},
		Dependencies: []string{
			"github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1.EditingOptions", "github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1.GitHubRepositoryConfig", "github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1.LocalRepositoryConfig", "github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1.S3RepositoryConfig"},
	}
}

func schema_pkg_apis_provisioning_v0alpha1_RepositoryStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "The status of a Repository. This is expected never to be created by a kubectl call or similar, and is expected to rarely (if ever) be edited manually. As such, it is also a little less well structured than the spec, such as conditional-but-ever-present fields.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"currentGitCommit": {
						SchemaProps: spec.SchemaProps{
							Description: "The Git commit we're currently synced to. A non-empty value only matters if we use a git storage backend. Useful for no-clone Git clients and if cloning Git clients ever lose their clones.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
	}
}

func schema_pkg_apis_provisioning_v0alpha1_ResourceWrapper(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "This is a container type for any resource type",
				Type:        []string{"object"},
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
					"path": {
						SchemaProps: spec.SchemaProps{
							Description: "Path to the remote file",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"ref": {
						SchemaProps: spec.SchemaProps{
							Description: "The commit hash (if exists)",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"hash": {
						SchemaProps: spec.SchemaProps{
							Description: "The repo hash value",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"timestamp": {
						SchemaProps: spec.SchemaProps{
							Description: "The modified time in the remote file system",
							Ref:         ref("k8s.io/apimachinery/pkg/apis/meta/v1.Time"),
						},
					},
					"value": {
						SchemaProps: spec.SchemaProps{
							Description: "The raw resource body",
							Ref:         ref("github.com/grafana/grafana/pkg/apimachinery/apis/common/v0alpha1.Unstructured"),
						},
					},
				},
				Required: []string{"value"},
			},
		},
		Dependencies: []string{
			"github.com/grafana/grafana/pkg/apimachinery/apis/common/v0alpha1.Unstructured", "k8s.io/apimachinery/pkg/apis/meta/v1.Time"},
	}
}

func schema_pkg_apis_provisioning_v0alpha1_S3RepositoryConfig(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"region": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"bucket": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
				},
			},
		},
	}
}

func schema_pkg_apis_provisioning_v0alpha1_WebhookResponse(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Dummy object to return for webhooks",
				Type:        []string{"object"},
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
					"status": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
				},
			},
		},
	}
}
