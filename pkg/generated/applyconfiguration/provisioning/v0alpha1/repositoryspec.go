// SPDX-License-Identifier: AGPL-3.0-only

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v0alpha1

import (
	provisioningv0alpha1 "github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1"
)

// RepositorySpecApplyConfiguration represents a declarative configuration of the RepositorySpec type for use
// with apply.
type RepositorySpecApplyConfiguration struct {
	Title       *string                                   `json:"title,omitempty"`
	Description *string                                   `json:"description,omitempty"`
	Folder      *string                                   `json:"folder,omitempty"`
	Editing     *EditingOptionsApplyConfiguration         `json:"editing,omitempty"`
	Type        *provisioningv0alpha1.RepositoryType      `json:"type,omitempty"`
	Local       *LocalRepositoryConfigApplyConfiguration  `json:"local,omitempty"`
	S3          *S3RepositoryConfigApplyConfiguration     `json:"s3,omitempty"`
	GitHub      *GitHubRepositoryConfigApplyConfiguration `json:"github,omitempty"`
}

// RepositorySpecApplyConfiguration constructs a declarative configuration of the RepositorySpec type for use with
// apply.
func RepositorySpec() *RepositorySpecApplyConfiguration {
	return &RepositorySpecApplyConfiguration{}
}

// WithTitle sets the Title field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Title field is set to the value of the last call.
func (b *RepositorySpecApplyConfiguration) WithTitle(value string) *RepositorySpecApplyConfiguration {
	b.Title = &value
	return b
}

// WithDescription sets the Description field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Description field is set to the value of the last call.
func (b *RepositorySpecApplyConfiguration) WithDescription(value string) *RepositorySpecApplyConfiguration {
	b.Description = &value
	return b
}

// WithFolder sets the Folder field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Folder field is set to the value of the last call.
func (b *RepositorySpecApplyConfiguration) WithFolder(value string) *RepositorySpecApplyConfiguration {
	b.Folder = &value
	return b
}

// WithEditing sets the Editing field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Editing field is set to the value of the last call.
func (b *RepositorySpecApplyConfiguration) WithEditing(value *EditingOptionsApplyConfiguration) *RepositorySpecApplyConfiguration {
	b.Editing = value
	return b
}

// WithType sets the Type field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Type field is set to the value of the last call.
func (b *RepositorySpecApplyConfiguration) WithType(value provisioningv0alpha1.RepositoryType) *RepositorySpecApplyConfiguration {
	b.Type = &value
	return b
}

// WithLocal sets the Local field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Local field is set to the value of the last call.
func (b *RepositorySpecApplyConfiguration) WithLocal(value *LocalRepositoryConfigApplyConfiguration) *RepositorySpecApplyConfiguration {
	b.Local = value
	return b
}

// WithS3 sets the S3 field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the S3 field is set to the value of the last call.
func (b *RepositorySpecApplyConfiguration) WithS3(value *S3RepositoryConfigApplyConfiguration) *RepositorySpecApplyConfiguration {
	b.S3 = value
	return b
}

// WithGitHub sets the GitHub field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the GitHub field is set to the value of the last call.
func (b *RepositorySpecApplyConfiguration) WithGitHub(value *GitHubRepositoryConfigApplyConfiguration) *RepositorySpecApplyConfiguration {
	b.GitHub = value
	return b
}
