package reststorage

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	secretv0alpha1 "github.com/grafana/grafana/pkg/apis/secret/v0alpha1"
	"github.com/grafana/grafana/pkg/registry/apis/secret/contracts"
	"github.com/grafana/grafana/pkg/registry/apis/secret/xkube"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func NewFakeKeeperStore(latency time.Duration) contracts.KeeperStorage {
	return &fakeKeeperStorage{
		values:  make(map[string]map[string]secretv0alpha1.Keeper),
		latency: latency,
	}
}

type fakeKeeperStorage struct {
	values  map[string]map[string]secretv0alpha1.Keeper
	latency time.Duration
}

func (s *fakeKeeperStorage) Create(ctx context.Context, k *secretv0alpha1.Keeper) (*secretv0alpha1.Keeper, error) {
	v := *k
	v.SetUID(types.UID(uuid.NewString()))
	v.ObjectMeta.SetResourceVersion(strconv.FormatInt(metav1.Now().UnixMicro(), 10))
	ns, ok := s.values[k.Namespace]
	if !ok {
		ns = make(map[string]secretv0alpha1.Keeper)
		s.values[k.Namespace] = ns
	}
	time.AfterFunc(s.latency, func() {
		ns[k.Name] = v
	})

	return &v, nil
}

func (s *fakeKeeperStorage) Read(ctx context.Context, nn xkube.NameNamespace) (*secretv0alpha1.Keeper, error) {
	ns, ok := s.values[nn.Namespace.String()]
	if !ok {
		return nil, contracts.ErrSecureValueNotFound
	}
	v, ok := ns[nn.Name]
	if !ok {
		return nil, contracts.ErrSecureValueNotFound
	}

	return &v, nil
}

func (s *fakeKeeperStorage) Update(ctx context.Context, nk *secretv0alpha1.Keeper) (*secretv0alpha1.Keeper, error) {
	v := *nk
	v.ObjectMeta.SetResourceVersion(strconv.FormatInt(metav1.Now().UnixMicro(), 10))
	ns, ok := s.values[nk.Namespace]
	if !ok {
		return nil, contracts.ErrSecureValueNotFound
	}
	_, ok = ns[nk.Name]
	if !ok {
		return nil, contracts.ErrSecureValueNotFound
	}
	time.AfterFunc(s.latency, func() {
		ns[nk.Name] = v
	})

	return &v, nil
}

func (s *fakeKeeperStorage) Delete(ctx context.Context, nn xkube.NameNamespace) error {
	ns, ok := s.values[nn.Namespace.String()]
	if !ok {
		return contracts.ErrSecureValueNotFound
	}
	_, ok = ns[nn.Name]
	if !ok {
		return contracts.ErrSecureValueNotFound
	}
	time.AfterFunc(s.latency, func() {
		delete(ns, nn.Name)
	})

	return nil
}

func (s *fakeKeeperStorage) List(ctx context.Context, namespace xkube.Namespace, options *internalversion.ListOptions) (*secretv0alpha1.KeeperList, error) {
	ns, ok := s.values[namespace.String()]
	if !ok {
		s.values[namespace.String()] = ns
	}
	l := make([]secretv0alpha1.Keeper, len(ns))
	for _, v := range ns {
		l = append(l, v)
	}
	return &secretv0alpha1.KeeperList{
		Items: l,
	}, nil
}
