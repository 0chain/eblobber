package handler

import (
	"context"
	"errors"
	"net/http"
	"path/filepath"

	"github.com/0chain/gosdk/constants"
	"gorm.io/gorm"

	"github.com/0chain/blobber/code/go/0chain.net/blobbercore/allocation"
	"github.com/0chain/blobber/code/go/0chain.net/blobbercore/reference"
	"github.com/0chain/blobber/code/go/0chain.net/core/common"
)

// DeleteFileCommand command for deleting file
type DeleteFileCommand struct {
	existingFileRef  *reference.Ref
	changeProcessor  *allocation.DeleteFileChange
	allocationChange *allocation.AllocationChange
	path             string
	connectionID     string
}

func (cmd *DeleteFileCommand) GetExistingFileRef() *reference.Ref {
	return cmd.existingFileRef
}

func (cmd *DeleteFileCommand) GetPath() string {
	return cmd.path
}

// IsValidated validate request.
func (cmd *DeleteFileCommand) IsValidated(ctx context.Context, req *http.Request, allocationObj *allocation.Allocation, clientID string) error {
	if allocationObj.OwnerID != clientID && allocationObj.RepairerID != clientID {
		return common.NewError("invalid_operation", "Operation needs to be performed by the owner or the payer of the allocation")
	}

	path, ok := common.GetField(req, "path")
	if !ok {
		return common.NewError("invalid_parameters", "Invalid path")
	}

	if filepath.Clean(path) != path {
		return common.NewError("invalid_parameters", "Invalid path")
	}

	cmd.path = path

	connectionID, ok := common.GetField(req, "connection_id")
	if !ok {
		return common.NewError("invalid_parameters", "Invalid connection id passed")
	}
	cmd.connectionID = connectionID
	var err error
	lookUpHash := reference.GetReferenceLookup(allocationObj.ID, path)
	cmd.existingFileRef, err = reference.GetLimitedRefFieldsByLookupHashWith(ctx, allocationObj.ID, lookUpHash, []string{"path", "name", "type", "id", "size"})
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return common.ErrFileWasDeleted
		}
		return common.NewError("bad_db_operation", err.Error())
	}
	if cmd.existingFileRef.Type == reference.DIRECTORY {
		// check if directory is empty
		empty, err := reference.IsDirectoryEmpty(ctx, cmd.existingFileRef.ID)
		if err != nil {
			return common.NewError("bad_db_operation", err.Error())
		}
		if !empty {
			return common.NewError("invalid_operation", "Directory is not empty")
		}
	}
	cmd.existingFileRef.LookupHash = lookUpHash
	return nil
}

// UpdateChange add DeleteFileChange in db
func (cmd *DeleteFileCommand) UpdateChange(ctx context.Context) error {
	err := cmd.AddChange(ctx)
	if err == gorm.ErrDuplicatedKey {
		return nil
	}
	return err
}

func (cmd *DeleteFileCommand) AddChange(ctx context.Context) error {
	connectionInput, _ := cmd.changeProcessor.Marshal()
	cmd.allocationChange.Input = connectionInput
	return cmd.allocationChange.Create(ctx)
}

// ProcessContent flush file to FileStorage
func (cmd *DeleteFileCommand) ProcessContent(_ context.Context, allocationObj *allocation.Allocation) (allocation.UploadResult, error) {
	deleteSize := cmd.existingFileRef.Size
	connectionID := cmd.connectionID
	cmd.changeProcessor = &allocation.DeleteFileChange{ConnectionID: connectionID,
		AllocationID: allocationObj.ID, Name: cmd.existingFileRef.Name,
		LookupHash: cmd.existingFileRef.LookupHash, Path: cmd.existingFileRef.Path, Size: deleteSize, Type: cmd.existingFileRef.Type}

	result := allocation.UploadResult{}
	result.Filename = cmd.existingFileRef.Name
	result.Size = cmd.existingFileRef.Size
	result.UpdateChange = true

	cmd.allocationChange = &allocation.AllocationChange{}
	cmd.allocationChange.ConnectionID = connectionID
	cmd.allocationChange.Size = 0 - deleteSize
	cmd.allocationChange.Operation = constants.FileOperationDelete
	cmd.allocationChange.LookupHash = cmd.existingFileRef.LookupHash

	allocation.UpdateConnectionObjSize(connectionID, cmd.allocationChange.Size)

	return result, nil
}

// ProcessThumbnail no thumbnail should be processed for delete. A deffered delete command has been added on ProcessContent
func (cmd *DeleteFileCommand) ProcessThumbnail(allocationObj *allocation.Allocation) error {
	//DO NOTHING
	return nil
}

func (cmd *DeleteFileCommand) GetNumBlocks() int64 {
	return 0
}
