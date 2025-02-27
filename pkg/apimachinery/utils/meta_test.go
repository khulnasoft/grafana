package utils_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/grafana/grafana/pkg/apimachinery/utils"
)

type TestResource struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec Spec `json:"spec,omitempty"`

	// Read/write raw status
	Status Spec `json:"status,omitempty"`
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TestResource) DeepCopyInto(out *TestResource) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Playlist.
func (in *TestResource) DeepCopy() *TestResource {
	if in == nil {
		return nil
	}
	out := new(TestResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TestResource) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// Spec defines model for Spec.
type Spec struct {
	// Name of the object.
	Title string `json:"title"`
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Spec) DeepCopyInto(out *Spec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Spec.
func (in *Spec) DeepCopy() *Spec {
	if in == nil {
		return nil
	}
	out := new(Spec)
	in.DeepCopyInto(out)
	return out
}

type TestResource2 struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec Spec2 `json:"spec,omitempty"`

	// Exercise read/write pointer status
	Status *Spec `json:"status,omitempty"`
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TestResource2) DeepCopyInto(out *TestResource2) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Playlist.
func (in *TestResource2) DeepCopy() *TestResource2 {
	if in == nil {
		return nil
	}
	out := new(TestResource2)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TestResource2) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// Spec defines model for Spec.
type Spec2 struct{}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Spec2) DeepCopyInto(out *Spec2) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Spec.
func (in *Spec2) DeepCopy() *Spec2 {
	if in == nil {
		return nil
	}
	out := new(Spec2)
	in.DeepCopyInto(out)
	return out
}

func TestMetaAccessor(t *testing.T) {
	repoInfo := &utils.ResourceRepositoryInfo{
		Name: "test",
		Path: "a/b/c",
		Hash: "kkk",
	}

	t.Run("fails for non resource objects", func(t *testing.T) {
		_, err := utils.MetaAccessor("hello")
		require.Error(t, err)

		_, err = utils.MetaAccessor(unstructured.Unstructured{})
		require.Error(t, err) // Not a pointer!

		_, err = utils.MetaAccessor(&unstructured.Unstructured{})
		require.NoError(t, err) // Must be a pointer

		_, err = utils.MetaAccessor(&TestResource{
			Spec: Spec{
				Title: "HELLO",
			},
		})
		require.NoError(t, err) // Must be a pointer
	})

	t.Run("get and set grafana labels (unstructured)", func(t *testing.T) {
		res := &unstructured.Unstructured{
			Object: map[string]any{},
		}
		meta, err := utils.MetaAccessor(res)
		require.NoError(t, err)

		// should return 0 when not set
		require.Equal(t, meta.GetDeprecatedInternalID(), int64(0))

		// 0 is not allowed
		meta.SetDeprecatedInternalID(0)
		require.Equal(t, map[string]string(nil), res.GetLabels())

		// should be able to set and get
		meta.SetDeprecatedInternalID(1)
		require.Equal(t, map[string]string{
			"grafana.app/deprecatedInternalID": "1",
		}, res.GetLabels())
		require.Equal(t, meta.GetDeprecatedInternalID(), int64(1))
	})

	t.Run("get and set grafana metadata (unstructured)", func(t *testing.T) {
		// Error reading spec+status when missing
		res := &unstructured.Unstructured{
			Object: map[string]any{},
		}
		meta, err := utils.MetaAccessor(res)
		require.NoError(t, err)
		spec, err := meta.GetSpec()
		require.Error(t, err)
		require.Nil(t, spec)
		status, err := meta.GetStatus()
		require.Error(t, err)
		require.Nil(t, status)

		// Now set a spec and status
		res.Object = map[string]any{
			"spec": map[string]any{
				"hello": "world",
				"title": "Title",
			},
			"status": map[string]any{
				"sloth": "🦥",
			},
		}

		meta.SetRepositoryInfo(repoInfo)
		meta.SetFolder("folderUID")

		require.Equal(t, map[string]string{
			"grafana.app/repoName": "test",
			"grafana.app/repoPath": "a/b/c",
			"grafana.app/repoHash": "kkk",
			"grafana.app/folder":   "folderUID",
		}, res.GetAnnotations())

		meta.SetNamespace("aaa")
		meta.SetResourceVersionInt64(12345)
		require.Equal(t, "aaa", res.GetNamespace())
		require.Equal(t, "aaa", meta.GetNamespace())

		rv, err := meta.GetResourceVersionInt64()
		require.NoError(t, err)
		require.Equal(t, int64(12345), rv)
		require.Equal(t, "Title", meta.FindTitle(""))

		// Make sure access to spec works for Unstructured
		spec, err = meta.GetSpec()
		require.NoError(t, err)
		require.Equal(t, res.Object["spec"], spec)
		spec = &map[string]string{"a": "b"}
		err = meta.SetSpec(spec)
		require.NoError(t, err)
		spec, err = meta.GetSpec()
		require.NoError(t, err)
		require.Equal(t, res.Object["spec"], spec)

		// Make sure access to spec works for Unstructured
		status, err = meta.GetStatus()
		require.NoError(t, err)
		require.Equal(t, res.Object["status"], status)
		status = &map[string]string{"a": "b"}
		err = meta.SetStatus(status)
		require.NoError(t, err)
		status, err = meta.GetStatus()
		require.NoError(t, err)
		require.Equal(t, res.Object["status"], status)
	})

	t.Run("get and set grafana metadata (TestResource)", func(t *testing.T) {
		res := &TestResource{
			Spec: Spec{
				Title: "test",
			},
			// Status is empty, but not nil!
		}
		meta, err := utils.MetaAccessor(res)
		require.NoError(t, err)

		meta.SetRepositoryInfo(repoInfo)
		meta.SetFolder("folderUID")

		require.Equal(t, map[string]string{
			"grafana.app/repoName": "test",
			"grafana.app/repoPath": "a/b/c",
			"grafana.app/repoHash": "kkk",
			"grafana.app/folder":   "folderUID",
		}, res.GetAnnotations())

		meta.SetNamespace("aaa")
		meta.SetResourceVersionInt64(12345)
		require.Equal(t, "aaa", res.GetNamespace())
		require.Equal(t, "aaa", meta.GetNamespace())

		rv, err := meta.GetResourceVersionInt64()
		require.NoError(t, err)
		require.Equal(t, int64(12345), rv)

		// Make sure access to spec works for Unstructured
		spec, err := meta.GetSpec()
		require.NoError(t, err)
		require.Equal(t, res.Spec, spec)
		err = meta.SetSpec(Spec{Title: "t2"})
		require.NoError(t, err)
		spec, err = meta.GetSpec()
		require.NoError(t, err)
		require.Equal(t, res.Spec, spec)
		require.Equal(t, `{"title":"t2"}`, asJSON(spec, false))

		// Check read/write status
		status, err := meta.GetStatus()
		require.NoError(t, err)
		require.NotNil(t, status)
		err = meta.SetStatus(Spec{Title: "111"})
		require.NoError(t, err)
		status, err = meta.GetStatus()
		require.NoError(t, err)
		require.Equal(t, res.Status, status)
		require.Equal(t, "111", res.Status.Title)
		require.Equal(t, `{"title":"111"}`, asJSON(status, false))
	})

	t.Run("get and set grafana metadata (TestResource2)", func(t *testing.T) {
		res := &TestResource2{
			Spec:   Spec2{},
			Status: &Spec{Title: "X"},
		}
		meta, err := utils.MetaAccessor(res)
		require.NoError(t, err)

		meta.SetRepositoryInfo(repoInfo)
		meta.SetFolder("folderUID")

		require.Equal(t, map[string]string{
			"grafana.app/repoName": "test",
			"grafana.app/repoPath": "a/b/c",
			"grafana.app/repoHash": "kkk",
			"grafana.app/folder":   "folderUID",
		}, res.GetAnnotations())

		meta.SetNamespace("aaa")
		meta.SetResourceVersionInt64(12345)
		require.Equal(t, "aaa", res.GetNamespace())
		require.Equal(t, "aaa", meta.GetNamespace())

		rv, err := meta.GetResourceVersionInt64()
		require.NoError(t, err)
		require.Equal(t, int64(12345), rv)

		// Make sure access to spec works for TestResource2
		spec, err := meta.GetSpec()
		require.NoError(t, err)
		require.Equal(t, res.Spec, spec)
		err = meta.SetSpec(Spec2{})
		require.NoError(t, err)
		spec, err = meta.GetSpec()
		require.NoError(t, err)
		require.Equal(t, res.Spec, spec)

		// Make sure access to spec works for TestResource2
		status, err := meta.GetStatus()
		require.NoError(t, err)
		require.Equal(t, res.Status, status)
		err = meta.SetStatus(&Spec{Title: "ZZ"})
		require.NoError(t, err)
		status, err = meta.GetStatus()
		require.NoError(t, err)
		require.Equal(t, res.Status, status)
		require.Equal(t, "ZZ", res.Status.Title)
	})

	t.Run("test reading old originInfo (now repository)", func(t *testing.T) {
		res := &TestResource2{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					"grafana.app/repoName": "test",
					"grafana.app/repoPath": "a/b/c",
					"grafana.app/repoHash": "zzz",
					"grafana.app/folder":   "folderUID",
				},
			},
			Spec: Spec2{},
		}
		meta, err := utils.MetaAccessor(res)
		require.NoError(t, err)

		info, err := meta.GetRepositoryInfo()
		require.NoError(t, err)
		require.Equal(t, "test", info.Name)
		require.Equal(t, "a/b/c", info.Path)
		require.Equal(t, "zzz", info.Hash)
	})

	t.Run("blob info", func(t *testing.T) {
		info := &utils.BlobInfo{UID: "AAA", Size: 123, Hash: "xyz", MimeType: "application/json", Charset: "utf-8"}
		anno := info.String()
		require.Equal(t, "AAA; size=123; hash=xyz; mime=application/json; charset=utf-8", anno)
		copy := utils.ParseBlobInfo(anno)
		require.Equal(t, info, copy)
	})

	t.Run("find titles", func(t *testing.T) {
		// with a k8s object that has Spec.Title
		obj := &TestResource{
			TypeMeta: metav1.TypeMeta{
				Kind:       "TestKIND",
				APIVersion: "aaa/v1alpha2",
			},
			Spec: Spec{
				Title: "HELLO",
			},
		}

		meta, err := utils.MetaAccessor(obj)
		require.NoError(t, err)
		meta.SetRepositoryInfo(repoInfo)
		meta.SetFolder("folderUID")

		require.Equal(t, map[string]string{
			"grafana.app/repoName": "test",
			"grafana.app/repoPath": "a/b/c",
			"grafana.app/repoHash": "kkk",
			"grafana.app/folder":   "folderUID",
		}, obj.GetAnnotations())

		require.Equal(t, "HELLO", obj.Spec.Title)
		require.Equal(t, "HELLO", meta.FindTitle(""))
		obj.Spec.Title = ""
		require.Equal(t, "", meta.FindTitle("xxx"))

		gvk := meta.GetGroupVersionKind()
		require.Equal(t, "aaa/v1alpha2, Kind=TestKIND", gvk.String())

		// with a k8s object without Spec.Title
		obj2 := &TestResource2{}

		meta, err = utils.MetaAccessor(obj2)
		require.NoError(t, err)
		meta.SetRepositoryInfo(repoInfo)
		meta.SetFolder("folderUID")

		require.Equal(t, map[string]string{
			"grafana.app/repoName": "test",
			"grafana.app/repoPath": "a/b/c",
			"grafana.app/repoHash": "kkk",
			"grafana.app/folder":   "folderUID",
		}, obj2.GetAnnotations())

		require.Equal(t, "xxx", meta.FindTitle("xxx"))

		rt, ok := meta.GetRuntimeObject()
		require.Equal(t, obj2, rt)
		require.True(t, ok)

		spec, err := meta.GetSpec()
		require.Equal(t, obj2.Spec, spec)
		require.NoError(t, err)
	})

	t.Run("ManagerProperties", func(t *testing.T) {
		tests := []struct {
			name           string
			setProperties  *utils.ManagerProperties
			wantProperties utils.ManagerProperties
			wantOK         bool
		}{
			{
				name: "get default values",
				wantProperties: utils.ManagerProperties{
					Identity:    "",
					Kind:        utils.ManagerKindUnknown,
					AllowsEdits: true,
					Suspended:   false,
				},
				wantOK: false,
			},
			{
				name: "set and get valid values",
				setProperties: &utils.ManagerProperties{
					Identity:    "identity",
					Kind:        utils.ManagerKindTerraform,
					AllowsEdits: false,
					Suspended:   false,
				},
				wantProperties: utils.ManagerProperties{
					Identity:    "identity",
					Kind:        utils.ManagerKindTerraform,
					AllowsEdits: false,
					Suspended:   false,
				},
				wantOK: true,
			},
			{
				name: "set empty identity returns default values",
				setProperties: &utils.ManagerProperties{
					Identity:    "",
					Kind:        utils.ManagerKindRepo,
					AllowsEdits: false,
					Suspended:   false,
				},
				wantProperties: utils.ManagerProperties{
					Identity:    "",
					Kind:        utils.ManagerKindUnknown,
					AllowsEdits: true,
					Suspended:   false,
				},
				wantOK: false,
			},
			{
				name: "invalid kind falls back to generic kind",
				setProperties: &utils.ManagerProperties{
					Identity:    "identity",
					Kind:        utils.ManagerKind("invalid"),
					AllowsEdits: false,
					Suspended:   true,
				},
				wantProperties: utils.ManagerProperties{
					Identity:    "identity",
					Kind:        utils.ManagerKindUnknown,
					AllowsEdits: false,
					Suspended:   true,
				},
				wantOK: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				res := &TestResource2{}
				meta, err := utils.MetaAccessor(res)
				require.NoError(t, err)

				if tt.setProperties != nil {
					meta.SetManagerProperties(*tt.setProperties)
				}

				mp, ok := meta.GetManagerProperties()
				require.Equal(t, tt.wantOK, ok)
				require.Equal(t, tt.wantProperties, mp)
			})
		}
	})

	t.Run("SourceProperties", func(t *testing.T) {
		tests := []struct {
			name           string
			setProperties  *utils.SourceProperties
			wantProperties utils.SourceProperties
			wantOK         bool
		}{
			{
				name:           "get default values",
				wantProperties: utils.SourceProperties{},
				wantOK:         false,
			},
			{
				name: "set and get valid values",
				setProperties: &utils.SourceProperties{
					Path:      "path",
					Checksum:  "hash",
					Timestamp: time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
				},
				wantProperties: utils.SourceProperties{
					Path:      "path",
					Checksum:  "hash",
					Timestamp: time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
				},
				wantOK: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				res := &TestResource2{}
				meta, err := utils.MetaAccessor(res)
				require.NoError(t, err)

				if tt.setProperties != nil {
					meta.SetSourceProperties(*tt.setProperties)
				}

				sp, ok := meta.GetSourceProperties()
				require.Equal(t, tt.wantProperties, sp)
				require.Equal(t, tt.wantOK, ok)
			})
		}
	})
}

func asJSON(v any, pretty bool) string {
	if v == nil {
		return ""
	}
	if pretty {
		bytes, _ := json.MarshalIndent(v, "", "  ")
		return string(bytes)
	}
	bytes, _ := json.Marshal(v)
	return string(bytes)
}
