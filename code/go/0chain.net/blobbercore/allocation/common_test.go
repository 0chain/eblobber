package allocation

import (
	"bytes"
	"io"
	"mime/multipart"

	"github.com/0chain/blobber/code/go/0chain.net/blobbercore/filestore"
)

type MockFileStore struct {
}

func (mfs *MockFileStore) Initialize() error {
	return nil
}

func (mfs *MockFileStore) WriteFile(allocID, connID string,
	fileData *filestore.FileInputData, infile multipart.File) (*filestore.FileOutputData, error) {

	b := bytes.NewBuffer(make([]byte, 0))
	n, _ := io.Copy(b, infile)
	return &filestore.FileOutputData{
		Name:            fileData.Name,
		Path:            fileData.Path,
		FixedMerkleRoot: "",
		ValidationRoot:  fileData.ValidationRoot,
		Size:            n,
	}, nil
}

func (mfs *MockFileStore) WriteRollback(allocID string, fileData *filestore.FileInputData) (*filestore.FileOutputData, error) {
	return nil, nil
}

func (mfs *MockFileStore) CommitWrite(allocID, connID string, fileData *filestore.FileInputData) (bool, error) {
	return true, nil
}

func (mfs *MockFileStore) DeleteTempFile(allocID, connID string, fileData *filestore.FileInputData) error {
	return nil
}

func (mfs *MockFileStore) DeleteFile(allocID string, contentHash, path, name string) error {
	return nil
}

func (mfs *MockFileStore) GetFileBlock(rin *filestore.ReadBlockInput) (*filestore.FileDownloadResponse, error) {
	return nil, nil
}

func (mfs *MockFileStore) GetBlocksMerkleTreeForChallenge(cir *filestore.ChallengeReadBlockInput) (*filestore.ChallengeResponse, error) {
	return nil, nil
}

func (mfs *MockFileStore) GetTotalTempFileSizes() (s uint64) {
	return 0
}

func (mfs *MockFileStore) GetTempFilesSizeOfAllocation(allocID string) uint64 {
	return 0
}

func (mfs *MockFileStore) GetTotalCommittedFileSize() uint64 {
	return 0
}

func (mfs *MockFileStore) GetCommittedFileSizeOfAllocation(allocID string) uint64 {
	return 0
}

func (mfs *MockFileStore) GetTotalFilesSize() uint64 {
	return 0
}

func (mfs *MockFileStore) GetTotalFilesSizeOfAllocation(allocID string) uint64 {
	return 0
}

func (mfs *MockFileStore) IterateObjects(allocationID string, handler filestore.FileObjectHandler) error {
	return nil
}

func (mfs *MockFileStore) GetCurrentDiskCapacity() uint64 {
	return 0
}

func (mfs *MockFileStore) CalculateCurrentDiskCapacity() error {
	return nil
}

func (mfs *MockFileStore) GetPathForFile(allocID, contentHash string) (string, error) {
	return "", nil
}

func (mfs *MockFileStore) UpdateAllocationMetaData(m map[string]interface{}) error {
	return nil
}
