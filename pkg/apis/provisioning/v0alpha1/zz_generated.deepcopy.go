//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// SPDX-License-Identifier: AGPL-3.0-only

// Code generated by deepcopy-gen. DO NOT EDIT.

package v0alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Author) DeepCopyInto(out *Author) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Author.
func (in *Author) DeepCopy() *Author {
	if in == nil {
		return nil
	}
	out := new(Author)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExportJobOptions) DeepCopyInto(out *ExportJobOptions) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExportJobOptions.
func (in *ExportJobOptions) DeepCopy() *ExportJobOptions {
	if in == nil {
		return nil
	}
	out := new(ExportJobOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FileItem) DeepCopyInto(out *FileItem) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FileItem.
func (in *FileItem) DeepCopy() *FileItem {
	if in == nil {
		return nil
	}
	out := new(FileItem)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FileList) DeepCopyInto(out *FileList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FileItem, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FileList.
func (in *FileList) DeepCopy() *FileList {
	if in == nil {
		return nil
	}
	out := new(FileList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FileList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GitHubRepositoryConfig) DeepCopyInto(out *GitHubRepositoryConfig) {
	*out = *in
	if in.EncryptedToken != nil {
		in, out := &in.EncryptedToken, &out.EncryptedToken
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GitHubRepositoryConfig.
func (in *GitHubRepositoryConfig) DeepCopy() *GitHubRepositoryConfig {
	if in == nil {
		return nil
	}
	out := new(GitHubRepositoryConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HealthStatus) DeepCopyInto(out *HealthStatus) {
	*out = *in
	if in.Message != nil {
		in, out := &in.Message, &out.Message
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HealthStatus.
func (in *HealthStatus) DeepCopy() *HealthStatus {
	if in == nil {
		return nil
	}
	out := new(HealthStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HistoryItem) DeepCopyInto(out *HistoryItem) {
	*out = *in
	if in.Authors != nil {
		in, out := &in.Authors, &out.Authors
		*out = make([]Author, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HistoryItem.
func (in *HistoryItem) DeepCopy() *HistoryItem {
	if in == nil {
		return nil
	}
	out := new(HistoryItem)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HistoryList) DeepCopyInto(out *HistoryList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]HistoryItem, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HistoryList.
func (in *HistoryList) DeepCopy() *HistoryList {
	if in == nil {
		return nil
	}
	out := new(HistoryList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HistoryList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Job) DeepCopyInto(out *Job) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Job.
func (in *Job) DeepCopy() *Job {
	if in == nil {
		return nil
	}
	out := new(Job)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Job) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobList) DeepCopyInto(out *JobList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Job, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobList.
func (in *JobList) DeepCopy() *JobList {
	if in == nil {
		return nil
	}
	out := new(JobList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *JobList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobResourceSummary) DeepCopyInto(out *JobResourceSummary) {
	*out = *in
	if in.Errors != nil {
		in, out := &in.Errors, &out.Errors
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobResourceSummary.
func (in *JobResourceSummary) DeepCopy() *JobResourceSummary {
	if in == nil {
		return nil
	}
	out := new(JobResourceSummary)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobSpec) DeepCopyInto(out *JobSpec) {
	*out = *in
	if in.PullRequest != nil {
		in, out := &in.PullRequest, &out.PullRequest
		*out = new(PullRequestJobOptions)
		**out = **in
	}
	if in.Export != nil {
		in, out := &in.Export, &out.Export
		*out = new(ExportJobOptions)
		**out = **in
	}
	if in.Sync != nil {
		in, out := &in.Sync, &out.Sync
		*out = new(SyncJobOptions)
		**out = **in
	}
	if in.Migrate != nil {
		in, out := &in.Migrate, &out.Migrate
		*out = new(MigrateJobOptions)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobSpec.
func (in *JobSpec) DeepCopy() *JobSpec {
	if in == nil {
		return nil
	}
	out := new(JobSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobStatus) DeepCopyInto(out *JobStatus) {
	*out = *in
	if in.Errors != nil {
		in, out := &in.Errors, &out.Errors
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Summary != nil {
		in, out := &in.Summary, &out.Summary
		*out = make([]*JobResourceSummary, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(JobResourceSummary)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobStatus.
func (in *JobStatus) DeepCopy() *JobStatus {
	if in == nil {
		return nil
	}
	out := new(JobStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LocalRepositoryConfig) DeepCopyInto(out *LocalRepositoryConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LocalRepositoryConfig.
func (in *LocalRepositoryConfig) DeepCopy() *LocalRepositoryConfig {
	if in == nil {
		return nil
	}
	out := new(LocalRepositoryConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MigrateJobOptions) DeepCopyInto(out *MigrateJobOptions) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MigrateJobOptions.
func (in *MigrateJobOptions) DeepCopy() *MigrateJobOptions {
	if in == nil {
		return nil
	}
	out := new(MigrateJobOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PullRequestJobOptions) DeepCopyInto(out *PullRequestJobOptions) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PullRequestJobOptions.
func (in *PullRequestJobOptions) DeepCopy() *PullRequestJobOptions {
	if in == nil {
		return nil
	}
	out := new(PullRequestJobOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Repository) DeepCopyInto(out *Repository) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Repository.
func (in *Repository) DeepCopy() *Repository {
	if in == nil {
		return nil
	}
	out := new(Repository)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Repository) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RepositoryList) DeepCopyInto(out *RepositoryList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Repository, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RepositoryList.
func (in *RepositoryList) DeepCopy() *RepositoryList {
	if in == nil {
		return nil
	}
	out := new(RepositoryList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RepositoryList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RepositorySpec) DeepCopyInto(out *RepositorySpec) {
	*out = *in
	if in.Workflows != nil {
		in, out := &in.Workflows, &out.Workflows
		*out = make([]Workflow, len(*in))
		copy(*out, *in)
	}
	out.Sync = in.Sync
	if in.Local != nil {
		in, out := &in.Local, &out.Local
		*out = new(LocalRepositoryConfig)
		**out = **in
	}
	if in.GitHub != nil {
		in, out := &in.GitHub, &out.GitHub
		*out = new(GitHubRepositoryConfig)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RepositorySpec.
func (in *RepositorySpec) DeepCopy() *RepositorySpec {
	if in == nil {
		return nil
	}
	out := new(RepositorySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RepositoryStatus) DeepCopyInto(out *RepositoryStatus) {
	*out = *in
	in.Health.DeepCopyInto(&out.Health)
	in.Sync.DeepCopyInto(&out.Sync)
	if in.Stats != nil {
		in, out := &in.Stats, &out.Stats
		*out = make([]ResourceCount, len(*in))
		copy(*out, *in)
	}
	if in.Webhook != nil {
		in, out := &in.Webhook, &out.Webhook
		*out = new(WebhookStatus)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RepositoryStatus.
func (in *RepositoryStatus) DeepCopy() *RepositoryStatus {
	if in == nil {
		return nil
	}
	out := new(RepositoryStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RepositoryView) DeepCopyInto(out *RepositoryView) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RepositoryView.
func (in *RepositoryView) DeepCopy() *RepositoryView {
	if in == nil {
		return nil
	}
	out := new(RepositoryView)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RepositoryViewList) DeepCopyInto(out *RepositoryViewList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]RepositoryView, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RepositoryViewList.
func (in *RepositoryViewList) DeepCopy() *RepositoryViewList {
	if in == nil {
		return nil
	}
	out := new(RepositoryViewList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RepositoryViewList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceCount) DeepCopyInto(out *ResourceCount) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceCount.
func (in *ResourceCount) DeepCopy() *ResourceCount {
	if in == nil {
		return nil
	}
	out := new(ResourceCount)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceList) DeepCopyInto(out *ResourceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ResourceListItem, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceList.
func (in *ResourceList) DeepCopy() *ResourceList {
	if in == nil {
		return nil
	}
	out := new(ResourceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceListItem) DeepCopyInto(out *ResourceListItem) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceListItem.
func (in *ResourceListItem) DeepCopy() *ResourceListItem {
	if in == nil {
		return nil
	}
	out := new(ResourceListItem)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceObjects) DeepCopyInto(out *ResourceObjects) {
	*out = *in
	out.Type = in.Type
	in.File.DeepCopyInto(&out.File)
	in.Existing.DeepCopyInto(&out.Existing)
	in.DryRun.DeepCopyInto(&out.DryRun)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceObjects.
func (in *ResourceObjects) DeepCopy() *ResourceObjects {
	if in == nil {
		return nil
	}
	out := new(ResourceObjects)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceStats) DeepCopyInto(out *ResourceStats) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ResourceCount, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceStats.
func (in *ResourceStats) DeepCopy() *ResourceStats {
	if in == nil {
		return nil
	}
	out := new(ResourceStats)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceStats) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceType) DeepCopyInto(out *ResourceType) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceType.
func (in *ResourceType) DeepCopy() *ResourceType {
	if in == nil {
		return nil
	}
	out := new(ResourceType)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceWrapper) DeepCopyInto(out *ResourceWrapper) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	if in.Timestamp != nil {
		in, out := &in.Timestamp, &out.Timestamp
		*out = (*in).DeepCopy()
	}
	in.Resource.DeepCopyInto(&out.Resource)
	if in.Errors != nil {
		in, out := &in.Errors, &out.Errors
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceWrapper.
func (in *ResourceWrapper) DeepCopy() *ResourceWrapper {
	if in == nil {
		return nil
	}
	out := new(ResourceWrapper)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceWrapper) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SyncJobOptions) DeepCopyInto(out *SyncJobOptions) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SyncJobOptions.
func (in *SyncJobOptions) DeepCopy() *SyncJobOptions {
	if in == nil {
		return nil
	}
	out := new(SyncJobOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SyncOptions) DeepCopyInto(out *SyncOptions) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SyncOptions.
func (in *SyncOptions) DeepCopy() *SyncOptions {
	if in == nil {
		return nil
	}
	out := new(SyncOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SyncStatus) DeepCopyInto(out *SyncStatus) {
	*out = *in
	if in.Message != nil {
		in, out := &in.Message, &out.Message
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SyncStatus.
func (in *SyncStatus) DeepCopy() *SyncStatus {
	if in == nil {
		return nil
	}
	out := new(SyncStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TestResults) DeepCopyInto(out *TestResults) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	if in.Errors != nil {
		in, out := &in.Errors, &out.Errors
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Details != nil {
		in, out := &in.Details, &out.Details
		*out = (*in).DeepCopy()
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TestResults.
func (in *TestResults) DeepCopy() *TestResults {
	if in == nil {
		return nil
	}
	out := new(TestResults)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TestResults) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WebhookResponse) DeepCopyInto(out *WebhookResponse) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	if in.Job != nil {
		in, out := &in.Job, &out.Job
		*out = new(JobSpec)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WebhookResponse.
func (in *WebhookResponse) DeepCopy() *WebhookResponse {
	if in == nil {
		return nil
	}
	out := new(WebhookResponse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *WebhookResponse) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WebhookStatus) DeepCopyInto(out *WebhookStatus) {
	*out = *in
	if in.EncryptedSecret != nil {
		in, out := &in.EncryptedSecret, &out.EncryptedSecret
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	if in.SubscribedEvents != nil {
		in, out := &in.SubscribedEvents, &out.SubscribedEvents
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WebhookStatus.
func (in *WebhookStatus) DeepCopy() *WebhookStatus {
	if in == nil {
		return nil
	}
	out := new(WebhookStatus)
	in.DeepCopyInto(out)
	return out
}
