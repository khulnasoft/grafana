// Code generated by mockery v2.42.2. DO NOT EDIT.

package dashboards

import (
	context "context"

	folder "github.com/grafana/grafana/pkg/services/folder"
	mock "github.com/stretchr/testify/mock"

	quota "github.com/grafana/grafana/pkg/services/quota"

	time "time"
)

// FakeDashboardStore is an autogenerated mock type for the Store type
type FakeDashboardStore struct {
	mock.Mock
}

// Count provides a mock function with given fields: _a0, _a1
func (_m *FakeDashboardStore) Count(_a0 context.Context, _a1 *quota.ScopeParameters) (*quota.Map, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Count")
	}

	var r0 *quota.Map
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *quota.ScopeParameters) (*quota.Map, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *quota.ScopeParameters) *quota.Map); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*quota.Map)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *quota.ScopeParameters) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CountDashboardsInFolders provides a mock function with given fields: ctx, request
func (_m *FakeDashboardStore) CountDashboardsInFolders(ctx context.Context, request *CountDashboardsInFolderRequest) (int64, error) {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for CountDashboardsInFolders")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *CountDashboardsInFolderRequest) (int64, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *CountDashboardsInFolderRequest) int64); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *CountDashboardsInFolderRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteDashboard provides a mock function with given fields: ctx, cmd
func (_m *FakeDashboardStore) DeleteDashboard(ctx context.Context, cmd *DeleteDashboardCommand) error {
	ret := _m.Called(ctx, cmd)

	if len(ret) == 0 {
		panic("no return value specified for DeleteDashboard")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *DeleteDashboardCommand) error); ok {
		r0 = rf(ctx, cmd)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAllDashboards provides a mock function with given fields: ctx, orgID
func (_m *FakeDashboardStore) DeleteAllDashboards(ctx context.Context, orgID int64) error {
	ret := _m.Called(ctx, orgID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteDashboard")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, orgID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteDashboardsInFolders provides a mock function with given fields: ctx, request
func (_m *FakeDashboardStore) DeleteDashboardsInFolders(ctx context.Context, request *DeleteDashboardsInFolderRequest) error {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for DeleteDashboardsInFolders")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *DeleteDashboardsInFolderRequest) error); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteOrphanedProvisionedDashboards provides a mock function with given fields: ctx, cmd
func (_m *FakeDashboardStore) DeleteOrphanedProvisionedDashboards(ctx context.Context, cmd *DeleteOrphanedProvisionedDashboardsCommand) error {
	ret := _m.Called(ctx, cmd)

	if len(ret) == 0 {
		panic("no return value specified for DeleteOrphanedProvisionedDashboards")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *DeleteOrphanedProvisionedDashboardsCommand) error); ok {
		r0 = rf(ctx, cmd)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindDashboards provides a mock function with given fields: ctx, query
func (_m *FakeDashboardStore) FindDashboards(ctx context.Context, query *FindPersistedDashboardsQuery) ([]DashboardSearchProjection, error) {
	ret := _m.Called(ctx, query)

	if len(ret) == 0 {
		panic("no return value specified for FindDashboards")
	}

	var r0 []DashboardSearchProjection
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *FindPersistedDashboardsQuery) ([]DashboardSearchProjection, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *FindPersistedDashboardsQuery) []DashboardSearchProjection); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]DashboardSearchProjection)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *FindPersistedDashboardsQuery) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllDashboards provides a mock function with given fields: ctx
func (_m *FakeDashboardStore) GetAllDashboards(ctx context.Context) ([]*Dashboard, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllDashboards")
	}

	var r0 []*Dashboard
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*Dashboard, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*Dashboard); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*Dashboard)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllDashboardsByOrgId provides a mock function with given fields: ctx
func (_m *FakeDashboardStore) GetAllDashboardsByOrgId(ctx context.Context, orgID int64) ([]*Dashboard, error) {
	ret := _m.Called(ctx, orgID)

	if len(ret) == 0 {
		panic("no return value specified for GetAllDashboardsByOrgId")
	}

	var r0 []*Dashboard
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) ([]*Dashboard, error)); ok {
		return rf(ctx, orgID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) []*Dashboard); ok {
		r0 = rf(ctx, orgID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*Dashboard)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, orgID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDashboard provides a mock function with given fields: ctx, query
func (_m *FakeDashboardStore) GetDashboard(ctx context.Context, query *GetDashboardQuery) (*Dashboard, error) {
	ret := _m.Called(ctx, query)

	if len(ret) == 0 {
		panic("no return value specified for GetDashboard")
	}

	var r0 *Dashboard
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *GetDashboardQuery) (*Dashboard, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *GetDashboardQuery) *Dashboard); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Dashboard)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *GetDashboardQuery) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDashboardTags provides a mock function with given fields: ctx, query
func (_m *FakeDashboardStore) GetDashboardTags(ctx context.Context, query *GetDashboardTagsQuery) ([]*DashboardTagCloudItem, error) {
	ret := _m.Called(ctx, query)

	if len(ret) == 0 {
		panic("no return value specified for GetDashboardTags")
	}

	var r0 []*DashboardTagCloudItem
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *GetDashboardTagsQuery) ([]*DashboardTagCloudItem, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *GetDashboardTagsQuery) []*DashboardTagCloudItem); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*DashboardTagCloudItem)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *GetDashboardTagsQuery) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDashboardUIDByID provides a mock function with given fields: ctx, query
func (_m *FakeDashboardStore) GetDashboardUIDByID(ctx context.Context, query *GetDashboardRefByIDQuery) (*DashboardRef, error) {
	ret := _m.Called(ctx, query)

	if len(ret) == 0 {
		panic("no return value specified for GetDashboardUIDByID")
	}

	var r0 *DashboardRef
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *GetDashboardRefByIDQuery) (*DashboardRef, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *GetDashboardRefByIDQuery) *DashboardRef); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*DashboardRef)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *GetDashboardRefByIDQuery) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDashboards provides a mock function with given fields: ctx, query
func (_m *FakeDashboardStore) GetDashboards(ctx context.Context, query *GetDashboardsQuery) ([]*Dashboard, error) {
	ret := _m.Called(ctx, query)

	if len(ret) == 0 {
		panic("no return value specified for GetDashboards")
	}

	var r0 []*Dashboard
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *GetDashboardsQuery) ([]*Dashboard, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *GetDashboardsQuery) []*Dashboard); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*Dashboard)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *GetDashboardsQuery) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDashboardsByPluginID provides a mock function with given fields: ctx, query
func (_m *FakeDashboardStore) GetDashboardsByPluginID(ctx context.Context, query *GetDashboardsByPluginIDQuery) ([]*Dashboard, error) {
	ret := _m.Called(ctx, query)

	if len(ret) == 0 {
		panic("no return value specified for GetDashboardsByPluginID")
	}

	var r0 []*Dashboard
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *GetDashboardsByPluginIDQuery) ([]*Dashboard, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *GetDashboardsByPluginIDQuery) []*Dashboard); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*Dashboard)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *GetDashboardsByPluginIDQuery) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProvisionedDashboardData provides a mock function with given fields: ctx, name
func (_m *FakeDashboardStore) GetProvisionedDashboardData(ctx context.Context, name string) ([]*DashboardProvisioning, error) {
	ret := _m.Called(ctx, name)

	if len(ret) == 0 {
		panic("no return value specified for GetProvisionedDashboardData")
	}

	var r0 []*DashboardProvisioning
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]*DashboardProvisioning, error)); ok {
		return rf(ctx, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []*DashboardProvisioning); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*DashboardProvisioning)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProvisionedDataByDashboardID provides a mock function with given fields: ctx, dashboardID
func (_m *FakeDashboardStore) GetProvisionedDataByDashboardID(ctx context.Context, dashboardID int64) (*DashboardProvisioning, error) {
	ret := _m.Called(ctx, dashboardID)

	if len(ret) == 0 {
		panic("no return value specified for GetProvisionedDataByDashboardID")
	}

	var r0 *DashboardProvisioning
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*DashboardProvisioning, error)); ok {
		return rf(ctx, dashboardID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *DashboardProvisioning); ok {
		r0 = rf(ctx, dashboardID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*DashboardProvisioning)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, dashboardID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProvisionedDataByDashboardUID provides a mock function with given fields: ctx, orgID, dashboardUID
func (_m *FakeDashboardStore) GetProvisionedDataByDashboardUID(ctx context.Context, orgID int64, dashboardUID string) (*DashboardProvisioning, error) {
	ret := _m.Called(ctx, orgID, dashboardUID)

	if len(ret) == 0 {
		panic("no return value specified for GetProvisionedDataByDashboardUID")
	}

	var r0 *DashboardProvisioning
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, string) (*DashboardProvisioning, error)); ok {
		return rf(ctx, orgID, dashboardUID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, string) *DashboardProvisioning); ok {
		r0 = rf(ctx, orgID, dashboardUID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*DashboardProvisioning)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, string) error); ok {
		r1 = rf(ctx, orgID, dashboardUID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSoftDeletedDashboard provides a mock function with given fields: ctx, orgID, uid
func (_m *FakeDashboardStore) GetSoftDeletedDashboard(ctx context.Context, orgID int64, uid string) (*Dashboard, error) {
	ret := _m.Called(ctx, orgID, uid)

	if len(ret) == 0 {
		panic("no return value specified for GetSoftDeletedDashboard")
	}

	var r0 *Dashboard
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, string) (*Dashboard, error)); ok {
		return rf(ctx, orgID, uid)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, string) *Dashboard); ok {
		r0 = rf(ctx, orgID, uid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Dashboard)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, string) error); ok {
		r1 = rf(ctx, orgID, uid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSoftDeletedExpiredDashboards provides a mock function with given fields: ctx, duration
func (_m *FakeDashboardStore) GetSoftDeletedExpiredDashboards(ctx context.Context, duration time.Duration) ([]*Dashboard, error) {
	ret := _m.Called(ctx, duration)

	if len(ret) == 0 {
		panic("no return value specified for GetSoftDeletedExpiredDashboards")
	}

	var r0 []*Dashboard
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, time.Duration) ([]*Dashboard, error)); ok {
		return rf(ctx, duration)
	}
	if rf, ok := ret.Get(0).(func(context.Context, time.Duration) []*Dashboard); ok {
		r0 = rf(ctx, duration)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*Dashboard)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, time.Duration) error); ok {
		r1 = rf(ctx, duration)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RestoreDashboard provides a mock function with given fields: ctx, orgID, dashboardUid, _a3
func (_m *FakeDashboardStore) RestoreDashboard(ctx context.Context, orgID int64, dashboardUid string, _a3 *folder.Folder) error {
	ret := _m.Called(ctx, orgID, dashboardUid, _a3)

	if len(ret) == 0 {
		panic("no return value specified for RestoreDashboard")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, string, *folder.Folder) error); ok {
		r0 = rf(ctx, orgID, dashboardUid, _a3)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveDashboard provides a mock function with given fields: ctx, cmd
func (_m *FakeDashboardStore) SaveDashboard(ctx context.Context, cmd SaveDashboardCommand) (*Dashboard, error) {
	ret := _m.Called(ctx, cmd)

	if len(ret) == 0 {
		panic("no return value specified for SaveDashboard")
	}

	var r0 *Dashboard
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, SaveDashboardCommand) (*Dashboard, error)); ok {
		return rf(ctx, cmd)
	}
	if rf, ok := ret.Get(0).(func(context.Context, SaveDashboardCommand) *Dashboard); ok {
		r0 = rf(ctx, cmd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Dashboard)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, SaveDashboardCommand) error); ok {
		r1 = rf(ctx, cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveProvisionedDashboard provides a mock function with given fields: ctx, cmd, provisioning
func (_m *FakeDashboardStore) SaveProvisionedDashboard(ctx context.Context, cmd SaveDashboardCommand, provisioning *DashboardProvisioning) (*Dashboard, error) {
	ret := _m.Called(ctx, cmd, provisioning)

	if len(ret) == 0 {
		panic("no return value specified for SaveProvisionedDashboard")
	}

	var r0 *Dashboard
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, SaveDashboardCommand, *DashboardProvisioning) (*Dashboard, error)); ok {
		return rf(ctx, cmd, provisioning)
	}
	if rf, ok := ret.Get(0).(func(context.Context, SaveDashboardCommand, *DashboardProvisioning) *Dashboard); ok {
		r0 = rf(ctx, cmd, provisioning)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Dashboard)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, SaveDashboardCommand, *DashboardProvisioning) error); ok {
		r1 = rf(ctx, cmd, provisioning)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SoftDeleteDashboard provides a mock function with given fields: ctx, orgID, dashboardUid
func (_m *FakeDashboardStore) SoftDeleteDashboard(ctx context.Context, orgID int64, dashboardUid string) error {
	ret := _m.Called(ctx, orgID, dashboardUid)

	if len(ret) == 0 {
		panic("no return value specified for SoftDeleteDashboard")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, string) error); ok {
		r0 = rf(ctx, orgID, dashboardUid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SoftDeleteDashboardsInFolders provides a mock function with given fields: ctx, orgID, folderUids
func (_m *FakeDashboardStore) SoftDeleteDashboardsInFolders(ctx context.Context, orgID int64, folderUids []string) error {
	ret := _m.Called(ctx, orgID, folderUids)

	if len(ret) == 0 {
		panic("no return value specified for SoftDeleteDashboardsInFolders")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, []string) error); ok {
		r0 = rf(ctx, orgID, folderUids)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UnprovisionDashboard provides a mock function with given fields: ctx, id
func (_m *FakeDashboardStore) UnprovisionDashboard(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for UnprovisionDashboard")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValidateDashboardBeforeSave provides a mock function with given fields: ctx, dashboard, overwrite
func (_m *FakeDashboardStore) ValidateDashboardBeforeSave(ctx context.Context, dashboard *Dashboard, overwrite bool) (bool, error) {
	ret := _m.Called(ctx, dashboard, overwrite)

	if len(ret) == 0 {
		panic("no return value specified for ValidateDashboardBeforeSave")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *Dashboard, bool) (bool, error)); ok {
		return rf(ctx, dashboard, overwrite)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *Dashboard, bool) bool); ok {
		r0 = rf(ctx, dashboard, overwrite)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *Dashboard, bool) error); ok {
		r1 = rf(ctx, dashboard, overwrite)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewFakeDashboardStore creates a new instance of FakeDashboardStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFakeDashboardStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *FakeDashboardStore {
	mock := &FakeDashboardStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
