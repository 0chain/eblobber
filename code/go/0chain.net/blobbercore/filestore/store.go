package filestore

import (
	"mime/multipart"
)

const (
	CHUNK_SIZE = 64 * 1024
	VERSION    = 1
)

type FileInputData struct {
	Name          string
	Path          string
	DataHash      string
	ThumbnailHash string
	LookupHash    string
	// ChunkSize chunk size
	ChunkSize int64
	//UploadLength indicates the size of the entire upload in bytes. The value MUST be a non-negative integer.
	UploadLength int64
	//Upload-Offset indicates a byte offset within a resource. The value MUST be a non-negative integer.
	UploadOffset int64
	//IsFinal  the request is final chunk
	IsFinal     bool
	IsThumbnail bool
	Size        int64
	Hasher      *CommitHasher
}

type FileOutputData struct {
	Name          string
	Path          string
	ThumbnailHash string
	// Size written size/chunk size
	Size int64
	// ChunkUploaded the chunk is uploaded or not.
	ChunkUploaded bool
	ContentSize   int64
}

type FileObjectHandler func(contentHash string, contentSize int64)

type FileStorer interface {
	// WriteFile write chunk file into disk
	Initialize() error
	WriteFile(allocID, connID string, fileData *FileInputData, infile multipart.File) (*FileOutputData, error)
	WriteDataToTree(allocID, connID, fileName, filePathHash string, hahser *CommitHasher) error
	CommitWrite(allocID, connID string, fileData *FileInputData) (bool, error)
	DeleteTempFile(allocID, connID string, fileData *FileInputData) error
	DeleteFile(allocID, contentHash string, version int) error
	MoveToFilestore(allocID, hash string, version int) error
	DeleteFromFilestore(allocID, hash string, version int) error
	DeletePreCommitDir(allocID string) error
	DeleteAllocation(allocID string)
	CopyFile(allocationID, oldFileLookupHash, newFileLookupHash string) error
	// GetFileBlock Get blocks of file starting from blockNum upto numBlocks. blockNum can't be less than 1.
	GetFileBlock(readBlockIn *ReadBlockInput) (*FileDownloadResponse, error)
	GetBlocksMerkleTreeForChallenge(cri *ChallengeReadBlockInput) (*ChallengeResponse, error)
	GetTotalTempFileSizes() (s uint64)
	GetTempFilesSizeOfAllocation(allocID string) uint64
	GetTotalCommittedFileSize() uint64
	GetCommittedFileSizeOfAllocation(allocID string) uint64
	GetFilePathSize(allocID, filehash, thumbHash string, version int) (int64, int64, error)
	GetTotalFilesSize() uint64
	GetTotalFilesSizeOfAllocation(allocID string) uint64
	GetTempFilePath(allocID, connID, fileName, filePathHash string) string

	IterateObjects(allocationID string, handler FileObjectHandler) error
	// SetupAllocation(allocationID string, skipCreate bool) (*StoreAllocation, error)
	GetCurrentDiskCapacity() uint64
	CalculateCurrentDiskCapacity() error
	// GetPathForFile given allocation id and content hash of a file, its path is calculated.
	// Will return error if allocation id or content hash are not of length 64
	GetPathForFile(allocID, contentHash string, version int) (string, error)
	// UpdateAllocationMetaData only updates if allocation size has changed or new allocation is allocated. Must use allocationID.
	// Use of allocation Tx might leak memory. allocation size must be of int64 type otherwise it won't be updated
	UpdateAllocationMetaData(m map[string]interface{}) error
}

var fileStore FileStorer

func SetFileStore(fs FileStorer) {
	fileStore = fs
}
func GetFileStore() FileStorer {
	return fileStore
}

// swagger:model FileDownloadResponse
type FileDownloadResponse struct {
	Nodes   [][][]byte
	Indexes [][]int
	Data    []byte
}

type ReadBlockInput struct {
	AllocationID     string
	FileSize         int64
	Hash             string
	StartBlockNum    int
	NumBlocks        int
	IsThumbnail      bool
	VerifyDownload   bool
	IsPrecommit      bool
	FilestoreVersion int
}

type ChallengeResponse struct {
	Proof   [][]byte `json:"proof"`
	Data    []byte   `json:"data"`
	LeafInd int      `json:"leaf_ind"`
}

type ChallengeReadBlockInput struct {
	BlockOffset      int
	FileSize         int64
	Hash             string
	AllocationID     string
	IsPrecommit      bool
	FilestoreVersion int
}
