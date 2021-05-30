// Code generated by mockery 2.7.5. DO NOT EDIT.

package mocks

import (
	io "io"

	filestore "github.com/0chain/blobber/code/go/0chain.net/blobbercore/filestore"

	json "encoding/json"

	mock "github.com/stretchr/testify/mock"

	multipart "mime/multipart"

	util "github.com/0chain/blobber/code/go/0chain.net/core/util"
)

// FileStore is an autogenerated mock type for the FileStore type
type FileStore struct {
	mock.Mock
}

// CommitWrite provides a mock function with given fields: allocationID, fileData, connectionID
func (_m *FileStore) CommitWrite(allocationID string, fileData *filestore.FileInputData, connectionID string) (bool, error) {
	ret := _m.Called(allocationID, fileData, connectionID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, *filestore.FileInputData, string) bool); ok {
		r0 = rf(allocationID, fileData, connectionID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *filestore.FileInputData, string) error); ok {
		r1 = rf(allocationID, fileData, connectionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteFile provides a mock function with given fields: allocationID, contentHash
func (_m *FileStore) DeleteFile(allocationID string, contentHash string) error {
	ret := _m.Called(allocationID, contentHash)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(allocationID, contentHash)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteTempFile provides a mock function with given fields: allocationID, fileData, connectionID
func (_m *FileStore) DeleteTempFile(allocationID string, fileData *filestore.FileInputData, connectionID string) error {
	ret := _m.Called(allocationID, fileData, connectionID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *filestore.FileInputData, string) error); ok {
		r0 = rf(allocationID, fileData, connectionID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DownloadFromCloud provides a mock function with given fields: fileHash, filePath
func (_m *FileStore) DownloadFromCloud(fileHash string, filePath string) error {
	ret := _m.Called(fileHash, filePath)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(fileHash, filePath)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetFileBlock provides a mock function with given fields: allocationID, fileData, blockNum, numBlocks
func (_m *FileStore) GetFileBlock(allocationID string, fileData *filestore.FileInputData, blockNum int64, numBlocks int64) ([]byte, error) {
	ret := _m.Called(allocationID, fileData, blockNum, numBlocks)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string, *filestore.FileInputData, int64, int64) []byte); ok {
		r0 = rf(allocationID, fileData, blockNum, numBlocks)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *filestore.FileInputData, int64, int64) error); ok {
		r1 = rf(allocationID, fileData, blockNum, numBlocks)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFileBlockForChallenge provides a mock function with given fields: allocationID, fileData, blockoffset
func (_m *FileStore) GetFileBlockForChallenge(allocationID string, fileData *filestore.FileInputData, blockoffset int) (json.RawMessage, util.MerkleTreeI, error) {
	ret := _m.Called(allocationID, fileData, blockoffset)

	var r0 json.RawMessage
	if rf, ok := ret.Get(0).(func(string, *filestore.FileInputData, int) json.RawMessage); ok {
		r0 = rf(allocationID, fileData, blockoffset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(json.RawMessage)
		}
	}

	var r1 util.MerkleTreeI
	if rf, ok := ret.Get(1).(func(string, *filestore.FileInputData, int) util.MerkleTreeI); ok {
		r1 = rf(allocationID, fileData, blockoffset)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(util.MerkleTreeI)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, *filestore.FileInputData, int) error); ok {
		r2 = rf(allocationID, fileData, blockoffset)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetTempPathSize provides a mock function with given fields: allocationID
func (_m *FileStore) GetTempPathSize(allocationID string) (int64, error) {
	ret := _m.Called(allocationID)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string) int64); ok {
		r0 = rf(allocationID)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(allocationID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTotalDiskSizeUsed provides a mock function with given fields:
func (_m *FileStore) GetTotalDiskSizeUsed() (int64, error) {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetlDiskSizeUsed provides a mock function with given fields: allocationID
func (_m *FileStore) GetlDiskSizeUsed(allocationID string) (int64, error) {
	ret := _m.Called(allocationID)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string) int64); ok {
		r0 = rf(allocationID)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(allocationID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IterateObjects provides a mock function with given fields: allocationID, handler
func (_m *FileStore) IterateObjects(allocationID string, handler filestore.FileObjectHandler) error {
	ret := _m.Called(allocationID, handler)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, filestore.FileObjectHandler) error); ok {
		r0 = rf(allocationID, handler)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetupAllocation provides a mock function with given fields: allocationID, skipCreate
func (_m *FileStore) SetupAllocation(allocationID string, skipCreate bool) (*filestore.StoreAllocation, error) {
	ret := _m.Called(allocationID, skipCreate)

	var r0 *filestore.StoreAllocation
	if rf, ok := ret.Get(0).(func(string, bool) *filestore.StoreAllocation); ok {
		r0 = rf(allocationID, skipCreate)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*filestore.StoreAllocation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, bool) error); ok {
		r1 = rf(allocationID, skipCreate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UploadToCloud provides a mock function with given fields: fileHash, filePath
func (_m *FileStore) UploadToCloud(fileHash string, filePath string) error {
	ret := _m.Called(fileHash, filePath)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(fileHash, filePath)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WriteFile provides a mock function with given fields: allocationID, fileData, infile, connectionID
func (_m *FileStore) WriteFile(allocationID string, fileData *filestore.FileInputData, infile multipart.File, connectionID string) (*filestore.FileOutputData, error) {
	ret := _m.Called(allocationID, fileData, infile, connectionID)

	var r0 *filestore.FileOutputData
	if rf, ok := ret.Get(0).(func(string, *filestore.FileInputData, multipart.File, string) *filestore.FileOutputData); ok {
		r0 = rf(allocationID, fileData, infile, connectionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*filestore.FileOutputData)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *filestore.FileInputData, multipart.File, string) error); ok {
		r1 = rf(allocationID, fileData, infile, connectionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WriteFileGRPC provides a mock function with given fields: allocationID, fileData, fileReader, connectionID
func (_m *FileStore) WriteFileGRPC(allocationID string, fileData *filestore.FileInputData, fileReader io.Reader, connectionID string) (*filestore.FileOutputData, error) {
	ret := _m.Called(allocationID, fileData, fileReader, connectionID)

	var r0 *filestore.FileOutputData
	if rf, ok := ret.Get(0).(func(string, *filestore.FileInputData, io.Reader, string) *filestore.FileOutputData); ok {
		r0 = rf(allocationID, fileData, fileReader, connectionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*filestore.FileOutputData)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *filestore.FileInputData, io.Reader, string) error); ok {
		r1 = rf(allocationID, fileData, fileReader, connectionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
