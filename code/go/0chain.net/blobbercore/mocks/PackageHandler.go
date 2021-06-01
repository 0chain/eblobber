// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	allocation "github.com/0chain/blobber/code/go/0chain.net/blobbercore/allocation"

	mock "github.com/stretchr/testify/mock"

	reference "github.com/0chain/blobber/code/go/0chain.net/blobbercore/reference"

	stats "github.com/0chain/blobber/code/go/0chain.net/blobbercore/stats"

	writemarker "github.com/0chain/blobber/code/go/0chain.net/blobbercore/writemarker"
)

// PackageHandler is an autogenerated mock type for the PackageHandler type
type PackageHandler struct {
	mock.Mock
}

// ApplyChanges provides a mock function with given fields: connectionObj, ctx, allocationRoot
func (_m *PackageHandler) ApplyChanges(connectionObj *allocation.AllocationChangeCollector, ctx context.Context, allocationRoot string) error {
	ret := _m.Called(connectionObj, ctx, allocationRoot)

	var r0 error
	if rf, ok := ret.Get(0).(func(*allocation.AllocationChangeCollector, context.Context, string) error); ok {
		r0 = rf(connectionObj, ctx, allocationRoot)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllocationChanges provides a mock function with given fields: ctx, connectionID, allocationID, clientID
func (_m *PackageHandler) GetAllocationChanges(ctx context.Context, connectionID string, allocationID string, clientID string) (*allocation.AllocationChangeCollector, error) {
	ret := _m.Called(ctx, connectionID, allocationID, clientID)

	var r0 *allocation.AllocationChangeCollector
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) *allocation.AllocationChangeCollector); ok {
		r0 = rf(ctx, connectionID, allocationID, clientID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*allocation.AllocationChangeCollector)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, connectionID, allocationID, clientID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCollaborators provides a mock function with given fields: ctx, refID
func (_m *PackageHandler) GetCollaborators(ctx context.Context, refID int64) ([]reference.Collaborator, error) {
	ret := _m.Called(ctx, refID)

	var r0 []reference.Collaborator
	if rf, ok := ret.Get(0).(func(context.Context, int64) []reference.Collaborator); ok {
		r0 = rf(ctx, refID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]reference.Collaborator)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, refID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCommitMetaTxns provides a mock function with given fields: ctx, refID
func (_m *PackageHandler) GetCommitMetaTxns(ctx context.Context, refID int64) ([]reference.CommitMetaTxn, error) {
	ret := _m.Called(ctx, refID)

	var r0 []reference.CommitMetaTxn
	if rf, ok := ret.Get(0).(func(context.Context, int64) []reference.CommitMetaTxn); ok {
		r0 = rf(ctx, refID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]reference.CommitMetaTxn)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, refID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFileStats provides a mock function with given fields: ctx, refID
func (_m *PackageHandler) GetFileStats(ctx context.Context, refID int64) (*stats.FileStats, error) {
	ret := _m.Called(ctx, refID)

	var r0 *stats.FileStats
	if rf, ok := ret.Get(0).(func(context.Context, int64) *stats.FileStats); ok {
		r0 = rf(ctx, refID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*stats.FileStats)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, refID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetObjectPath provides a mock function with given fields: ctx, allocationID, blockNum
func (_m *PackageHandler) GetObjectPath(ctx context.Context, allocationID string, blockNum int64) (*reference.ObjectPath, error) {
	ret := _m.Called(ctx, allocationID, blockNum)

	var r0 *reference.ObjectPath
	if rf, ok := ret.Get(0).(func(context.Context, string, int64) *reference.ObjectPath); ok {
		r0 = rf(ctx, allocationID, blockNum)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*reference.ObjectPath)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, int64) error); ok {
		r1 = rf(ctx, allocationID, blockNum)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetObjectTree provides a mock function with given fields: ctx, allocationID, path
func (_m *PackageHandler) GetObjectTree(ctx context.Context, allocationID string, path string) (*reference.Ref, error) {
	ret := _m.Called(ctx, allocationID, path)

	var r0 *reference.Ref
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *reference.Ref); ok {
		r0 = rf(ctx, allocationID, path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*reference.Ref)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, allocationID, path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRefWithChildren provides a mock function with given fields: ctx, allocationID, path
func (_m *PackageHandler) GetRefWithChildren(ctx context.Context, allocationID string, path string) (*reference.Ref, error) {
	ret := _m.Called(ctx, allocationID, path)

	var r0 *reference.Ref
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *reference.Ref); ok {
		r0 = rf(ctx, allocationID, path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*reference.Ref)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, allocationID, path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetReference provides a mock function with given fields: ctx, allocationID, path
func (_m *PackageHandler) GetReference(ctx context.Context, allocationID string, path string) (*reference.Ref, error) {
	ret := _m.Called(ctx, allocationID, path)

	var r0 *reference.Ref
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *reference.Ref); ok {
		r0 = rf(ctx, allocationID, path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*reference.Ref)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, allocationID, path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetReferenceFromLookupHash provides a mock function with given fields: ctx, allocationID, path_hash
func (_m *PackageHandler) GetReferenceFromLookupHash(ctx context.Context, allocationID string, path_hash string) (*reference.Ref, error) {
	ret := _m.Called(ctx, allocationID, path_hash)

	var r0 *reference.Ref
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *reference.Ref); ok {
		r0 = rf(ctx, allocationID, path_hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*reference.Ref)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, allocationID, path_hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetReferencePathFromPaths provides a mock function with given fields: ctx, allocationID, paths
func (_m *PackageHandler) GetReferencePathFromPaths(ctx context.Context, allocationID string, paths []string) (*reference.Ref, error) {
	ret := _m.Called(ctx, allocationID, paths)

	var r0 *reference.Ref
	if rf, ok := ret.Get(0).(func(context.Context, string, []string) *reference.Ref); ok {
		r0 = rf(ctx, allocationID, paths)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*reference.Ref)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, []string) error); ok {
		r1 = rf(ctx, allocationID, paths)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWriteMarkerEntity provides a mock function with given fields: ctx, allocation_root
func (_m *PackageHandler) GetWriteMarkerEntity(ctx context.Context, allocation_root string) (*writemarker.WriteMarkerEntity, error) {
	ret := _m.Called(ctx, allocation_root)

	var r0 *writemarker.WriteMarkerEntity
	if rf, ok := ret.Get(0).(func(context.Context, string) *writemarker.WriteMarkerEntity); ok {
		r0 = rf(ctx, allocation_root)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*writemarker.WriteMarkerEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, allocation_root)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsACollaborator provides a mock function with given fields: ctx, refID, clientID
func (_m *PackageHandler) IsACollaborator(ctx context.Context, refID int64, clientID string) bool {
	ret := _m.Called(ctx, refID, clientID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, int64, string) bool); ok {
		r0 = rf(ctx, refID, clientID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// UpdateAllocationAndCommitChanges provides a mock function with given fields: ctx, writemarkerObj, connectionObj, allocationObj, allocationRoot
func (_m *PackageHandler) UpdateAllocationAndCommitChanges(ctx context.Context, writemarkerObj *writemarker.WriteMarkerEntity, connectionObj *allocation.AllocationChangeCollector, allocationObj *allocation.Allocation, allocationRoot string) error {
	ret := _m.Called(ctx, writemarkerObj, connectionObj, allocationObj, allocationRoot)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *writemarker.WriteMarkerEntity, *allocation.AllocationChangeCollector, *allocation.Allocation, string) error); ok {
		r0 = rf(ctx, writemarkerObj, connectionObj, allocationObj, allocationRoot)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VerifyMarker provides a mock function with given fields: wm, ctx, sa, co
func (_m *PackageHandler) VerifyMarker(wm *writemarker.WriteMarkerEntity, ctx context.Context, sa *allocation.Allocation, co *allocation.AllocationChangeCollector) error {
	ret := _m.Called(wm, ctx, sa, co)

	var r0 error
	if rf, ok := ret.Get(0).(func(*writemarker.WriteMarkerEntity, context.Context, *allocation.Allocation, *allocation.AllocationChangeCollector) error); ok {
		r0 = rf(wm, ctx, sa, co)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
