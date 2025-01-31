package handler

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/0chain/blobber/code/go/0chain.net/blobbercore/blobberhttp"
	"github.com/0chain/blobber/code/go/0chain.net/blobbercore/config"

	"github.com/0chain/gosdk/constants"

	"github.com/0chain/blobber/code/go/0chain.net/blobbercore/allocation"
	"github.com/0chain/blobber/code/go/0chain.net/blobbercore/datastore"
	"github.com/0chain/blobber/code/go/0chain.net/blobbercore/filestore"
	"github.com/0chain/blobber/code/go/0chain.net/blobbercore/readmarker"
	"github.com/0chain/blobber/code/go/0chain.net/blobbercore/reference"
	"github.com/0chain/blobber/code/go/0chain.net/blobbercore/writemarker"
	"github.com/0chain/blobber/code/go/0chain.net/core/common"
	"github.com/0chain/blobber/code/go/0chain.net/core/encryption"
	"github.com/0chain/blobber/code/go/0chain.net/core/logging"
	"github.com/0chain/blobber/code/go/0chain.net/core/node"

	"go.uber.org/zap"
	"gorm.io/gorm"

	. "github.com/0chain/blobber/code/go/0chain.net/core/logging"
)

const (
	// EncryptionHeaderSize encryption header size in chunk: PRE.MessageChecksum(128)"+PRE.OverallChecksum(128)
	EncryptionHeaderSize = 128 + 128
	// ReEncryptionHeaderSize re-encryption header size in chunk
	ReEncryptionHeaderSize = 256
)

func readPreRedeem(
	ctx context.Context, alloc *allocation.Allocation,
	numBlocks, pendNumBlocks int64, payerID string) (err error) {

	if numBlocks == 0 {
		return
	}

	// check out read pool tokens if read_price > 0
	var (
		blobberID = node.Self.ID
	)

	if alloc.GetRequiredReadBalance(blobberID, numBlocks) <= 0 {
		return // skip if read price is zero
	}

	readPoolBalance, err := allocation.GetReadPoolsBalance(ctx, payerID)
	if err != nil {
		return common.NewError("read_pre_redeem", "database error while reading read pools balance")
	}

	requiredBalance := alloc.GetRequiredReadBalance(blobberID, numBlocks+pendNumBlocks)
	if float64(readPoolBalance) < requiredBalance {
		rp, err := allocation.RequestReadPoolStat(payerID)
		if err != nil {
			return common.NewErrorf("read_pre_redeem", "can't request read pools from sharders: %v", err)
		}

		rp.ClientID = payerID
		err = allocation.UpdateReadPool(ctx, rp)
		if err != nil {
			return common.NewErrorf("read_pre_redeem", "can't save requested read pools: %v", err)
		}

		readPoolBalance = rp.Balance

		if float64(readPoolBalance) < requiredBalance {
			return common.NewError("read_pre_redeem",
				"not enough tokens in client's read pools associated with the allocation->blobber")
		}
	}

	return
}

func writePreRedeem(ctx context.Context, alloc *allocation.Allocation, writeMarker *writemarker.WriteMarker, payerID string) (err error) {
	// check out read pool tokens if read_price > 0
	var (
		blobberID       = node.Self.ID
		requiredBalance = alloc.GetRequiredWriteBalance(blobberID, writeMarker.Size, writeMarker.Timestamp)
		wp              *allocation.WritePool
	)

	if writeMarker.Size <= 0 || requiredBalance <= 0 {
		return
	}

	writePoolBalance, err := allocation.GetWritePoolsBalance(ctx, alloc.ID)
	if err != nil {
		Logger.Error("write_pre_redeem:get_write_pools_balance", zap.Error(err), zap.String("allocation_id", alloc.ID))
		return common.NewError("write_pre_redeem", "database error while getting write pool balance")
	}

	pendingWriteSize, err := allocation.GetPendingWrite(ctx, payerID, alloc.ID)
	if err != nil {
		escapedPayerID := sanitizeString(payerID)
		Logger.Error("write_pre_redeem:get_pending_write", zap.Error(err), zap.String("allocation_id", alloc.ID), zap.String("payer_id", escapedPayerID))
		return common.NewError("write_pre_redeem", "database error while getting pending writes")
	}

	requiredBalance = alloc.GetRequiredWriteBalance(blobberID, pendingWriteSize+writeMarker.Size, writeMarker.Timestamp)

	if writePoolBalance < requiredBalance {
		wp, err = allocation.RequestWritePool(alloc.ID)
		if err != nil {
			return common.NewErrorf("write_pre_redeem", "can't request write pools from sharders: %v", err)
		}

		err = allocation.SetWritePool(ctx, alloc.ID, wp)
		if err != nil {
			return common.NewErrorf("write_pre_redeem", "can't save requested write pools: %v", err)
		}

		writePoolBalance += wp.Balance
	}

	if writePoolBalance < requiredBalance {
		return common.NewErrorf("write_pre_redeem", "not enough "+
			"tokens in write pools (client -> allocation ->  blobber)"+
			"(%s -> %s -> %s), available balance %d, required balance %d", payerID,
			alloc.ID, writeMarker.BlobberID, writePoolBalance, requiredBalance)
	}

	if err := allocation.AddToPending(ctx, payerID, alloc.ID, writeMarker.Size); err != nil {
		Logger.Error(err.Error())
		return common.NewErrorf("write_pre_redeem", "can't save pending writes in DB")

	}
	return
}

func (fsh *StorageHandler) RedeemReadMarker(ctx context.Context, r *http.Request) (interface{}, error) {
	var (
		clientID     = ctx.Value(constants.ContextKeyClient).(string)
		allocationTx = ctx.Value(constants.ContextKeyAllocation).(string)
		allocationID = ctx.Value(constants.ContextKeyAllocationID).(string)
		alloc        *allocation.Allocation
		blobberID    = node.Self.ID
		quotaManager = getQuotaManager()
	)

	if clientID == "" {
		return nil, common.NewError("redeem_readmarker", "invalid client")
	}

	alloc, err := fsh.verifyAllocation(ctx, allocationID, allocationTx, false)
	if err != nil {
		return nil, common.NewErrorf("redeem_readmarker", "invalid allocation id passed: %v", err)
	}

	dr, err := FromDownloadRequest(alloc.ID, r, true)
	if err != nil {
		return nil, err
	}

	isReadFree := alloc.IsReadFree(blobberID)
	if isReadFree {
		Logger.Info("free_read: readmarker not saved",
			zap.String("clientID", clientID),
			zap.String("allocationID", allocationID))
		return &blobberhttp.DownloadResponse{
			Success: true,
		}, nil
	}

	key := clientID + ":" + alloc.ID
	lock, isNewLock := readmarker.ReadmarkerMapLock.GetLock(key)
	if !isNewLock {
		return nil, common.NewErrorf("lock_exists", fmt.Sprintf("lock exists for key: %v", key))
	}

	lock.Lock()
	defer lock.Unlock()

	// create read marker
	var (
		rme              *readmarker.ReadMarkerEntity
		latestRM         *readmarker.ReadMarker
		latestRedeemedRC int64
		pendNumBlocks    int64
	)

	rme, err = readmarker.GetLatestReadMarkerEntity(ctx, clientID, alloc.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.NewErrorf("redeem_readmarker", "couldn't get read marker from DB: %v", err)
	}

	if rme != nil {
		latestRM = rme.LatestRM
		latestRedeemedRC = rme.LatestRedeemedRC
		if pendNumBlocks, err = rme.PendNumBlocks(); err != nil {
			return nil, common.NewErrorf("redeem_readmarker", "couldn't get number of blocks pending redeeming: %v", err)
		}
	}

	// check out read pool tokens if read_price > 0
	err = readPreRedeem(ctx, alloc, dr.ReadMarker.SessionRC, pendNumBlocks, clientID)
	if err != nil {
		return nil, common.NewErrorf("not_enough_tokens", "pre-redeeming read marker: %v", err)
	}

	if latestRM != nil && latestRM.ReadCounter+(dr.ReadMarker.SessionRC) > dr.ReadMarker.ReadCounter {
		latestRM.BlobberID = node.Self.ID
		return &blobberhttp.DownloadResponse{
			Success:  false,
			LatestRM: latestRM,
		}, common.NewError("stale_read_marker", "")
	}

	if dr.ReadMarker.ClientID != clientID {
		return nil, common.NewError("invalid_client", "header clientID and readmarker clientID are different")
	}

	rmObj := new(readmarker.ReadMarkerEntity)
	rmObj.LatestRM = &dr.ReadMarker

	if err = rmObj.VerifyMarker(ctx, alloc); err != nil {
		return nil, common.NewErrorf("redeem_readmarker", "invalid read marker, "+"failed to verify the read marker: %v", err)
	}

	err = readmarker.SaveLatestReadMarker(ctx, &dr.ReadMarker, latestRedeemedRC, latestRM == nil)
	if err != nil {
		Logger.Error(err.Error())
		return nil, common.NewError("redeem_readmarker", "couldn't save latest read marker")
	}

	quotaManager.createOrUpdateQuota(dr.ReadMarker.SessionRC, dr.ConnectionID)
	Logger.Info("readmarker_saved", zap.Any("rmObj", rmObj))
	return &blobberhttp.DownloadResponse{
		Success:  true,
		LatestRM: &dr.ReadMarker,
	}, nil
}

// swagger:route GET /v1/file/download/{allocation} GetDownloadFile
// Download a file.
//
// Download Handler (downloadFile). The response is either a byte stream or a FileDownloadResponse, which contains the file data or the thumbnail data, and the merkle proof if the download is verified.
// This depends on the "X-Verify-Download" header. If the header is set to "true", the response is a FileDownloadResponse, otherwise it is a byte stream.
//
// parameters:
//   +name: allocation
//     description: TxHash of the allocation in question.
//     in: path
//     required: true
//     type: string
//	 +name: X-App-Client-ID
//     description: The ID/Wallet address of the client sending the request.
//     in: header
//     type: string
//     required: true
//	 +name: X-App-Client-Key
// 	   description: The key of the client sending the request.
//     in: header
//     type: string
//     required: true
//	 +name: ALLOCATION-ID
//	   description: The ID of the allocation in question.
//     in: header
//     type: string
//     required: true
//  +name: X-Connection-ID
//     description: The ID of the connection used for the download. Usually, the download process occurs in multiple requests, on per block, where all of them are done in a single connection between the client and the blobber.
//	   in: header
//     type: string
//     required: false
//  +name: X-Path-Hash
//     description: The hash of the path of the file to download. If not provided, will be calculated from "X-Path" parameter.
//     in: header
//     type: string
//	   required: false
//  +name: X-Path
//     description: The path of the file to download.
//     in: header
//     type: string
//     required: true
//  +name: X-Block-Num
//     description: The block number of the file to download. Must be 0 or greater (valid index).
//     in: header
//     type: integer
//     required: false
//     default: 0
//  +name: X-Num-Blocks
//     description: The number of blocks to download. Must be 0 or greater.
//     in: header
//     type: integer
//     required: false
//     default: 0
//  +name: X-Read-Marker
//     description: The read marker to use for the download (check [ReadMarker](#/responses/ReadMarker)).
//     in: header
//     type: string
//     required: false
//  +name: X-Auth-Token
//     description: The auth token to use for the download. If the file is shared, the auth token is required.
//     in: header
//     type: string
//  +name: X-Mode
//     description: Download mode. Either "full" for full file download, or "thumbnail" to download the thumbnail of the file
//     in: header
//     type: string
//  +name: X-Verify-Download
//     description: If set to "true", the download should be verified. If the mode is "thumbnail", the thumbnail hash stored in the db is compared with the hash of the actual file. If the mode is "full", merkle proof is calculated and returned in the response.
//     in: header
//     type: string
//  +name: X-Version
//     description: If its value is "v2" then both allocation_id and blobber url base are hashed and verified using X-App-Client-Signature-V2.
//     in: header
//     type: string
//  +name: X-App-Client-Signature
//     description: Digital signature of the client used to verify the request if the X-Version is not "v2"
//     in: header
//     type: string
//  +name: X-App-Client-Signature-V2
//     description: Digital signature of the client used to verify the request if the X-Version is "v2"
//     in: header
//     type: string
//
// responses:
//
//   200: FileDownloadResponse
//   400:

func (fsh *StorageHandler) DownloadFile(ctx context.Context, r *http.Request) (interface{}, error) {
	// get client and allocation ids

	var (
		clientID        = ctx.Value(constants.ContextKeyClient).(string)
		clientPublicKey = ctx.Value(constants.ContextKeyClientKey).(string)
		allocationTx    = ctx.Value(constants.ContextKeyAllocation).(string)
		allocationID    = ctx.Value(constants.ContextKeyAllocationID).(string)
		alloc           *allocation.Allocation
		blobberID       = node.Self.ID
		quotaManager    = getQuotaManager()
	)

	if clientID == "" || clientPublicKey == "" {
		return nil, common.NewError("download_file", "invalid client")
	}

	if ok := CheckBlacklist(clientID); ok {
		return nil, common.NewError("blacklisted_client", "Client is blacklisted: "+clientID)
	}

	alloc, err := fsh.verifyAllocation(ctx, allocationID, allocationTx, false)
	if err != nil {
		return nil, common.NewErrorf("download_file", "invalid allocation id passed: %v", err)
	}

	dr, err := FromDownloadRequest(alloc.ID, r, false)
	if err != nil {
		return nil, err
	}

	if dr.NumBlocks > config.Configuration.BlockLimitRequest {
		return nil, common.NewErrorf("download_file", "too many blocks requested: %v, max limit is %v", dr.NumBlocks, config.Configuration.BlockLimitRequest)
	}

	dailyBlocksConsumed := getDailyBlocks(clientID)
	if dailyBlocksConsumed+dr.NumBlocks > config.Configuration.BlockLimitDaily {
		return nil, common.NewErrorf("download_file", "daily block limit reached: %v, max limit is %v", dailyBlocksConsumed, config.Configuration.BlockLimitDaily)
	}

	fileref, err := reference.GetReferenceByLookupHash(ctx, alloc.ID, dr.PathHash)
	if err != nil {
		return nil, common.NewErrorf("download_file", "invalid file path: %v", err)
	}

	if fileref.Type != reference.FILE {
		return nil, common.NewErrorf("download_file", "path is not a file: %v", err)
	}

	isOwner := clientID == alloc.OwnerID

	var authToken *readmarker.AuthTicket
	var shareInfo *reference.ShareInfo

	if !isOwner {
		if dr.AuthToken == "" {
			return nil, common.NewError("invalid_authticket", "authticket is required")
		}
		if dr.Version == "v2" {
			valid, err := verifySignatureFromRequest(allocationTx, r.Header.Get(common.ClientSignatureHeader), r.Header.Get(common.ClientSignatureHeaderV2), clientPublicKey)
			if !valid || err != nil {
				return nil, common.NewError("invalid_signature", "Invalid signature")
			}
		}
		authTokenString, err := base64.StdEncoding.DecodeString(dr.AuthToken)
		if err != nil {
			return nil, common.NewError("invalid_authticket", err.Error())
		}

		if authToken, err = fsh.verifyAuthTicket(ctx, string(authTokenString), alloc, fileref, clientID, false); authToken == nil {
			return nil, common.NewErrorf("invalid_authticket", "cannot verify auth ticket: %v", err)
		}

		shareInfo, err = reference.GetShareInfo(ctx, authToken.ClientID, authToken.FilePathHash)
		if err != nil || shareInfo == nil {
			return nil, common.NewError("invalid_share", "client does not have permission to download the file. share does not exist")
		}

		if shareInfo.Revoked {
			return nil, common.NewError("invalid_share", "client does not have permission to download the file. share revoked")
		}

		availableAt := shareInfo.AvailableAt.Unix()
		if common.Timestamp(availableAt) > common.Now() {
			return nil, common.NewErrorf("download_file", "the file is not available until: %v", shareInfo.AvailableAt.UTC().Format("2006-01-02T15:04:05"))
		}

	} else {
		if dr.Version == "v2" {
			valid, err := verifySignatureFromRequest(allocationTx, r.Header.Get(common.ClientSignatureHeader), r.Header.Get(common.ClientSignatureHeaderV2), alloc.OwnerPublicKey)
			if !valid || err != nil {
				return nil, common.NewError("invalid_signature", "Invalid signature")
			}
		}
	}

	isReadFree := alloc.IsReadFree(blobberID)
	var dq *DownloadQuota

	if !isReadFree {
		dq = quotaManager.getDownloadQuota(dr.ConnectionID)
		if dq == nil {
			return nil, common.NewError("download_file", fmt.Sprintf("no download quota for %v", dr.ConnectionID))
		}
		if dq.Quota < dr.NumBlocks {
			return nil, common.NewError("download_file", fmt.Sprintf("insufficient quota: available %v, requested %v", dq.Quota, dr.NumBlocks))
		}
	}

	var (
		downloadMode         = dr.DownloadMode
		fileDownloadResponse *filestore.FileDownloadResponse
	)

	if dr.BlockNum > math.MaxInt32 || dr.NumBlocks > math.MaxInt32 {
		return nil, common.NewErrorf("download_file", "BlockNum or NumBlocks is too large to convert to int")
	}

	fromPreCommit := alloc.AllocationVersion == fileref.AllocationVersion
	if downloadMode == DownloadContentThumb {
		rbi := &filestore.ReadBlockInput{
			AllocationID:     alloc.ID,
			FileSize:         fileref.ThumbnailSize,
			Hash:             fileref.LookupHash,
			StartBlockNum:    int(dr.BlockNum),
			NumBlocks:        int(dr.NumBlocks),
			IsThumbnail:      true,
			IsPrecommit:      fromPreCommit,
			FilestoreVersion: filestore.VERSION,
		}

		logging.Logger.Info("calling GetFileBlock for thumb", zap.Any("rbi", rbi))
		fileDownloadResponse, err = filestore.GetFileStore().GetFileBlock(rbi)
		if err != nil {
			return nil, common.NewErrorf("download_file", "couldn't get thumbnail block: %v", err)
		}
	} else {
		rbi := &filestore.ReadBlockInput{
			AllocationID:     alloc.ID,
			FileSize:         fileref.Size,
			Hash:             fileref.LookupHash,
			StartBlockNum:    int(dr.BlockNum),
			NumBlocks:        int(dr.NumBlocks),
			VerifyDownload:   dr.VerifyDownload,
			IsPrecommit:      fromPreCommit,
			FilestoreVersion: filestore.VERSION,
		}
		logging.Logger.Info("calling GetFileBlock", zap.Any("rbi", rbi))
		fileDownloadResponse, err = filestore.GetFileStore().GetFileBlock(rbi)
		if err != nil {
			return nil, common.NewErrorf("download_file", "couldn't get file block: %v", err)
		}
	}

	var chunkEncoder ChunkEncoder
	if len(fileref.EncryptedKey) > 0 && authToken != nil {
		chunkEncoder = &PREChunkEncoder{
			EncryptedKey:              fileref.EncryptedKey,
			ReEncryptionKey:           shareInfo.ReEncryptionKey,
			ClientEncryptionPublicKey: shareInfo.ClientEncryptionPublicKey,
		}
	} else {
		chunkEncoder = &RawChunkEncoder{}
	}

	chunkData, err := chunkEncoder.Encode(int(fileref.ChunkSize), fileDownloadResponse.Data)
	if err != nil {
		return nil, err
	}

	if !isReadFree {
		err = quotaManager.consumeQuota(dr.ConnectionID, dr.NumBlocks)
		if err != nil {
			return nil, common.NewError("download_file", err.Error())
		}
	}

	fileDownloadResponse.Data = chunkData
	reference.FileBlockDownloaded(ctx, fileref, dr.NumBlocks)
	go func() {
		addDailyBlocks(clientID, dr.NumBlocks)
	}()
	if !dr.VerifyDownload {
		return fileDownloadResponse.Data, nil
	}
	return fileDownloadResponse, nil
}

func (fsh *StorageHandler) CreateConnection(ctx context.Context, r *http.Request) (interface{}, error) {
	allocationTx := ctx.Value(constants.ContextKeyAllocation).(string)
	allocationId := ctx.Value(constants.ContextKeyAllocationID).(string)
	allocationObj, err := fsh.verifyAllocation(ctx, allocationId, allocationTx, false)
	if err != nil {
		return nil, common.NewError("invalid_parameters", "Invalid allocation id passed."+err.Error())
	}

	clientID := ctx.Value(constants.ContextKeyClient).(string)
	_ = ctx.Value(constants.ContextKeyClientKey).(string)

	if clientID == "" {
		return nil, common.NewError("invalid_operation", "Invalid client")
	}

	if allocationObj.OwnerID != clientID && allocationObj.RepairerID != clientID {
		return nil, common.NewError("invalid_operation", "Operation needs to be performed by the owner or the payer of the allocation")
	}

	valid, err := verifySignatureFromRequest(allocationTx, r.Header.Get(common.ClientSignatureHeader), r.Header.Get(common.ClientSignatureHeaderV2), allocationObj.OwnerPublicKey)
	if !valid || err != nil {
		return nil, common.NewError("invalid_signature", "Invalid signature")
	}

	connectionID := r.FormValue("connection_id")
	if connectionID == "" {
		return nil, common.NewError("invalid_parameters", "Invalid connection id passed")
	}

	connectionObj, err := allocation.GetAllocationChanges(ctx, connectionID, allocationObj.ID, clientID)
	if err != nil {
		return nil, common.NewError("meta_error", "Error reading metadata for connection")
	}
	err = connectionObj.Save(ctx)
	if err != nil {
		Logger.Error("Error in writing the connection meta data", zap.Error(err))
		return nil, common.NewError("connection_write_error", "Error writing the connection meta data")
	}

	return &blobberhttp.ConnectionResult{
		ConnectionID:   connectionID,
		AllocationRoot: allocationObj.AllocationRoot,
	}, nil
}

func (fsh *StorageHandler) CommitWrite(ctx context.Context, r *http.Request) (*blobberhttp.CommitResult, error) {
	startTime := time.Now()
	if r.Method == "GET" {
		return nil, common.NewError("invalid_method", "Invalid method used for the upload URL. Use POST instead")
	}

	allocationId := ctx.Value(constants.ContextKeyAllocationID).(string)
	allocationTx := ctx.Value(constants.ContextKeyAllocation).(string)
	clientID := ctx.Value(constants.ContextKeyClient).(string)
	clientKey := ctx.Value(constants.ContextKeyClientKey).(string)
	clientKeyBytes, _ := hex.DecodeString(clientKey)

	logging.Logger.Info("commit_write", zap.String("allocation_id", allocationId))

	if clientID == "" || clientKey == "" {
		return nil, common.NewError("invalid_parameters", "Please provide clientID and clientKey")
	}

	if ok := CheckBlacklist(clientID); ok {
		return nil, common.NewError("blacklisted_client", "Client is blacklisted: "+clientID)
	}

	allocationObj, err := fsh.verifyAllocation(ctx, allocationId, allocationTx, false)
	if err != nil {
		return nil, common.NewError("invalid_parameters", "Invalid allocation id passed."+err.Error())
	}
	if allocationObj.FileOptions == 0 {
		return nil, common.NewError("immutable_allocation", "Cannot write to an immutable allocation")
	}

	elapsedAllocation := time.Since(startTime)

	allocationID := allocationObj.ID

	connectionID, ok := common.GetField(r, "connection_id")
	if !ok {
		return nil, common.NewError("invalid_parameters", "Invalid connection id passed")
	}

	elapsedGetLock := time.Since(startTime) - elapsedAllocation

	connectionObj, err := allocation.GetAllocationChanges(ctx, connectionID, allocationID, clientID)
	if err != nil {
		// might be good to check if blobber already has stored writemarker
		return nil, common.NewErrorf("invalid_parameters",
			"Invalid connection id. Connection id was not found: %v", err)
	}
	if len(connectionObj.Changes) == 0 {
		logging.Logger.Info("commit_write_empty", zap.String("connection_id", connectionID))
	}

	if len(connectionObj.Changes) > config.Configuration.MaxConnectionChanges {
		return nil, common.NewError("max_connection_changes",
			"Max connection changes reached. A connection can only have "+fmt.Sprintf("%v", config.Configuration.MaxConnectionChanges)+" changes")
	}

	elapsedGetConnObj := time.Since(startTime) - elapsedAllocation - elapsedGetLock

	if allocationObj.OwnerID != clientID || encryption.Hash(clientKeyBytes) != clientID {
		return nil, common.NewError("invalid_operation", "Operation needs to be performed by the owner of the allocation")
	}

	if allocationObj.BlobberSizeUsed+connectionObj.Size > allocationObj.BlobberSize {
		return nil, common.NewError("max_allocation_size",
			"Max size reached for the allocation with this blobber")
	}

	var result blobberhttp.CommitResult
	versionMarkerStr := r.FormValue("version_marker")
	if versionMarkerStr == "" {
		return nil, common.NewError("invalid_parameters", "Invalid version marker passed")
	}
	versionMarker := writemarker.VersionMarker{}
	err = json.Unmarshal([]byte(versionMarkerStr), &versionMarker)
	if err != nil {
		return nil, common.NewError("unmarshall_error", fmt.Sprintf("Error while unmarshalling version marker: %s", err.Error()))
	}

	err = versionMarker.Verify(allocationID, allocationObj.OwnerPublicKey)
	if err != nil {
		return nil, err
	}

	// Move preCommitDir to finalDir
	if allocationObj.IsRedeemRequired {
		err = connectionObj.MoveToFilestore(ctx, allocationObj.AllocationVersion)
		if err != nil {
			return nil, common.NewError("move_to_filestore_error", fmt.Sprintf("Error while moving to filestore: %s", err.Error()))
		}
	}

	elapsedMoveToFilestore := time.Since(startTime) - elapsedAllocation - elapsedGetLock - elapsedGetConnObj
	if len(connectionObj.Changes) > 0 {
		err = connectionObj.ApplyChanges(
			ctx, common.Now(), versionMarker.Version)
		if err != nil {
			Logger.Error("Error applying changes", zap.Error(err))
			return nil, err
		}
	}

	elapsedApplyChanges := time.Since(startTime) - elapsedAllocation - elapsedGetLock -
		elapsedGetConnObj - elapsedMoveToFilestore

	db := datastore.GetStore().GetTransaction(ctx)
	err = db.Create(&versionMarker).Error
	if err != nil {
		return nil, common.NewError("db_error", fmt.Sprintf("Error while saving version marker: %s", err.Error()))
	}
	allocationObj.PrevBlobberSizeUsed = allocationObj.BlobberSizeUsed
	allocationObj.PrevUsedSize = allocationObj.UsedSize
	allocationObj.BlobberSizeUsed += connectionObj.Size
	allocationObj.UsedSize += connectionObj.Size
	allocationObj.AllocationVersion = versionMarker.Version
	if len(connectionObj.Changes) == 0 {
		allocationObj.IsRedeemRequired = false
	} else {
		allocationObj.IsRedeemRequired = true
	}

	updateMap := map[string]interface{}{
		"used_size":              allocationObj.UsedSize,
		"blobber_size_used":      allocationObj.BlobberSizeUsed,
		"is_redeem_required":     allocationObj.IsRedeemRequired,
		"allocation_version":     versionMarker.Version,
		"prev_used_size":         allocationObj.PrevUsedSize,
		"prev_blobber_size_used": allocationObj.PrevBlobberSizeUsed,
	}
	updateOption := func(a *allocation.Allocation) {
		a.IsRedeemRequired = allocationObj.IsRedeemRequired
		a.BlobberSizeUsed = allocationObj.BlobberSizeUsed
		a.UsedSize = allocationObj.UsedSize
		a.AllocationVersion = allocationObj.AllocationVersion
		a.PrevUsedSize = allocationObj.PrevUsedSize
		a.PrevBlobberSizeUsed = allocationObj.PrevBlobberSizeUsed
	}

	if err = allocation.Repo.UpdateAllocation(ctx, allocationObj, updateMap, updateOption); err != nil {
		return nil, common.NewError("allocation_write_error", "Error persisting the allocation object "+err.Error())
	}

	elapsedSaveAllocation := time.Since(startTime) - elapsedAllocation - elapsedGetLock -
		elapsedGetConnObj - elapsedApplyChanges

	err = connectionObj.CommitToFileStore(ctx)
	if err != nil {
		if !errors.Is(common.ErrFileWasDeleted, err) {
			return nil, common.NewError("file_store_error", "Error committing to file store. "+err.Error())
		}
	}
	elapsedCommitStore := time.Since(startTime) - elapsedAllocation - elapsedGetLock - elapsedGetConnObj - elapsedApplyChanges - elapsedSaveAllocation
	logging.Logger.Info("commit_filestore", zap.String("allocation_id", allocationId))
	connectionObj.DeleteChanges(ctx)

	db.Model(connectionObj).Updates(allocation.AllocationChangeCollector{Status: allocation.CommittedConnection})
	result.AllocationRoot = allocationObj.AllocationRoot
	result.Success = true
	result.ErrorMessage = ""
	var (
		commitOperation string
		input           string
	)
	if len(connectionObj.Changes) > 0 {
		commitOperation = connectionObj.Changes[0].Operation
		input = connectionObj.Changes[0].Input
	} else {
		commitOperation = "[commit]empty"
		input = "[commit]empty"
	}

	//Delete connection object and its changes

	db.Delete(connectionObj)
	go allocation.DeleteConnectionObjEntry(connectionID)
	go AddWriteMarkerCount(clientID, connectionObj.Size <= 0)

	Logger.Info("[commit]"+commitOperation,
		zap.String("alloc_id", allocationID),
		zap.String("input", input),
		zap.Duration("get_alloc", elapsedAllocation),
		zap.Duration("get-lock", elapsedGetLock),
		zap.Duration("get-conn-obj", elapsedGetConnObj),
		zap.Duration("move-to-filestore", elapsedMoveToFilestore),
		zap.Duration("apply-changes", elapsedApplyChanges),
		zap.Duration("save-allocation", elapsedSaveAllocation),
		zap.Duration("commit-store", elapsedCommitStore),
		zap.Duration("total", time.Since(startTime)),
	)
	return &result, nil
}

func (fsh *StorageHandler) RenameObject(ctx context.Context, r *http.Request) (interface{}, error) {
	allocationTx := ctx.Value(constants.ContextKeyAllocation).(string)
	allocationId := ctx.Value(constants.ContextKeyAllocationID).(string)
	allocationObj, err := fsh.verifyAllocation(ctx, allocationId, allocationTx, false)
	if err != nil {
		return nil, common.NewError("invalid_parameters", "Invalid allocation id passed."+err.Error())
	}

	if !allocationObj.CanRename() {
		return nil, common.NewError("prohibited_allocation_file_options", "Cannot rename data in this allocation.")
	}

	allocationID := allocationObj.ID

	clientID := ctx.Value(constants.ContextKeyClient).(string)
	_ = ctx.Value(constants.ContextKeyClientKey).(string)
	valid, err := verifySignatureFromRequest(allocationTx, r.Header.Get(common.ClientSignatureHeader), r.Header.Get(common.ClientSignatureHeaderV2), allocationObj.OwnerPublicKey)
	if !valid || err != nil {
		return nil, common.NewError("invalid_signature", "Invalid signature")
	}

	if clientID == "" {
		return nil, common.NewError("invalid_operation", "Invalid client")
	}

	new_name := r.FormValue("new_name")
	if new_name == "" {
		return nil, common.NewError("invalid_parameters", "Invalid name")
	}
	if filepath.Base(new_name) != new_name {
		logging.Logger.Error("invalid_parameters", zap.String("new_name", new_name), zap.String("base", filepath.Base(new_name)))
		return nil, common.NewError("invalid_parameters", "Invalid name")
	}
	pathHash, err := pathHashFromReq(r, allocationID)
	if err != nil {
		return nil, err
	}

	if clientID == "" || allocationObj.OwnerID != clientID {
		return nil, common.NewError("invalid_operation", "Operation needs to be performed by the owner of the allocation")
	}

	connectionID := r.FormValue("connection_id")
	if connectionID == "" {
		return nil, common.NewError("invalid_parameters", "Invalid connection id passed")
	}

	connectionObj, err := allocation.GetAllocationChanges(ctx, connectionID, allocationID, clientID)
	if err != nil {
		return nil, common.NewError("meta_error", "Error reading metadata for connection")
	}

	objectRef, err := reference.GetLimitedRefFieldsByLookupHash(ctx, allocationID, pathHash, []string{"id", "name", "path", "size", "type"})

	if err != nil {
		return nil, common.NewError("invalid_parameters", "Invalid file path. "+err.Error())
	}

	if objectRef.Path == "/" {
		return nil, common.NewError("invalid_operation", "cannot rename root path")
	}

	if objectRef.Type != reference.FILE {
		isEmpty, err := reference.IsDirectoryEmpty(ctx, objectRef.ID)
		if err != nil {
			return nil, common.NewError("invalid_operation", "Error checking if directory is empty")
		}
		if !isEmpty {
			return nil, common.NewError("invalid_operation", "Directory is not empty")
		}
	}

	allocationChange := &allocation.AllocationChange{}
	allocationChange.ConnectionID = connectionObj.ID
	allocationChange.Size = 0
	allocationChange.LookupHash = pathHash
	allocationChange.Operation = constants.FileOperationRename
	dfc := &allocation.RenameFileChange{ConnectionID: connectionObj.ID,
		AllocationID: connectionObj.AllocationID, Path: objectRef.Path, Type: objectRef.Type}
	dfc.NewName = new_name
	connectionObj.AddChange(allocationChange, dfc)

	err = connectionObj.Save(ctx)
	if err != nil {
		Logger.Error("Error in writing the connection meta data", zap.Error(err))
		return nil, common.NewError("connection_write_error", "Error writing the connection meta data")
	}

	result := &allocation.UploadResult{}
	result.Filename = new_name
	result.Size = objectRef.Size

	return result, nil
}

func (fsh *StorageHandler) CopyObject(ctx context.Context, r *http.Request) (interface{}, error) {

	allocationTx := ctx.Value(constants.ContextKeyAllocation).(string)
	allocationId := ctx.Value(constants.ContextKeyAllocationID).(string)
	allocationObj, err := fsh.verifyAllocation(ctx, allocationId, allocationTx, false)
	if err != nil {
		return nil, common.NewError("invalid_parameters", "Invalid allocation id passed."+err.Error())
	}

	if !allocationObj.CanCopy() {
		return nil, common.NewError("prohibited_allocation_file_options", "Cannot copy data from this allocation.")
	}

	valid, err := verifySignatureFromRequest(allocationTx, r.Header.Get(common.ClientSignatureHeader), r.Header.Get(common.ClientSignatureHeaderV2), allocationObj.OwnerPublicKey)
	if !valid || err != nil {
		return nil, common.NewError("invalid_signature", "Invalid signature")
	}

	clientID := ctx.Value(constants.ContextKeyClient).(string)
	_ = ctx.Value(constants.ContextKeyClientKey).(string)

	allocationID := allocationObj.ID

	if clientID == "" {
		return nil, common.NewError("invalid_operation", "Invalid client")
	}

	destPath := r.FormValue("dest")
	if destPath == "" {
		return nil, common.NewError("invalid_parameters", "Invalid destination for operation")
	}
	destPath = filepath.Clean(destPath)
	pathHash, err := pathHashFromReq(r, allocationID)
	if err != nil {
		return nil, err
	}

	if clientID == "" || allocationObj.OwnerID != clientID {
		return nil, common.NewError("invalid_operation", "Operation needs to be performed by the owner of the allocation")
	}

	connectionID := r.FormValue("connection_id")
	if connectionID == "" {
		return nil, common.NewError("invalid_parameters", "Invalid connection id passed")
	}

	connectionObj, err := allocation.GetAllocationChanges(ctx, connectionID, allocationID, clientID)
	if err != nil {
		return nil, common.NewError("meta_error", "Error reading metadata for connection")
	}

	objectRef, err := reference.GetLimitedRefFieldsByLookupHash(ctx, allocationID, pathHash, []string{"id", "name", "path", "size", "type"})

	if err != nil {
		return nil, common.NewError("invalid_parameters", "Invalid file path. "+err.Error())
	}
	if objectRef.ParentPath == destPath || objectRef.Path == destPath {
		return nil, common.NewError("invalid_parameters", "Invalid destination path. Cannot copy to the same parent directory.")
	}
	if objectRef.Type == reference.DIRECTORY {
		isEmpty, err := reference.IsDirectoryEmpty(ctx, objectRef.ID)
		if err != nil {
			return nil, common.NewError("invalid_operation", "Error checking if directory is empty")
		}
		if !isEmpty {
			return nil, common.NewError("invalid_operation", "Directory is not empty")
		}
	}
	newPath := filepath.Join(destPath, objectRef.Name)
	newPath = filepath.Clean(newPath)
	paths, err := common.GetParentPaths(newPath)
	if err != nil {
		return nil, err
	}

	paths = append(paths, newPath)

	refs, err := reference.GetRefsTypeFromPaths(ctx, allocationID, paths)
	if err != nil {
		Logger.Error("Database error", zap.Error(err))
		return nil, common.NewError("database_error", fmt.Sprintf("Got db error while getting refs for %v", paths))
	}

	for _, ref := range refs {
		switch ref.Path {
		case newPath:
			return nil, common.NewError("invalid_parameters", "Invalid destination path. Object Already exists.")
		default:
			if ref.Type == reference.FILE {
				return nil, common.NewError("invalid_path", fmt.Sprintf("%v is of file type", ref.Path))
			}
		}
	}

	allocationChange := &allocation.AllocationChange{}
	allocationChange.ConnectionID = connectionObj.ID
	allocationChange.Size = objectRef.Size
	allocationChange.LookupHash = pathHash
	allocationChange.Operation = constants.FileOperationCopy
	dfc := &allocation.CopyFileChange{ConnectionID: connectionObj.ID,
		AllocationID: connectionObj.AllocationID, DestPath: newPath, Type: objectRef.Type, SrcPath: objectRef.Path}
	allocation.UpdateConnectionObjSize(connectionID, allocationChange.Size)
	connectionObj.AddChange(allocationChange, dfc)

	err = connectionObj.Save(ctx)
	if err != nil {
		Logger.Error("Error in writing the connection meta data", zap.Error(err))
		return nil, common.NewError("connection_write_error", "Error writing the connection meta data")
	}

	result := &allocation.UploadResult{}
	result.Filename = objectRef.Name
	result.Size = objectRef.Size
	return result, nil
}

func (fsh *StorageHandler) MoveObject(ctx context.Context, r *http.Request) (interface{}, error) {

	allocationId := ctx.Value(constants.ContextKeyAllocationID).(string)
	allocationTx := ctx.Value(constants.ContextKeyAllocation).(string)
	allocationObj, err := fsh.verifyAllocation(ctx, allocationId, allocationTx, false)
	if err != nil {
		return nil, common.NewError("invalid_parameters", "Invalid allocation id passed."+err.Error())
	}

	if !allocationObj.CanMove() {
		return nil, common.NewError("prohibited_allocation_file_options", "Cannot move data in this allocation.")
	}

	valid, err := verifySignatureFromRequest(allocationTx, r.Header.Get(common.ClientSignatureHeader), r.Header.Get(common.ClientSignatureHeaderV2), allocationObj.OwnerPublicKey)
	if !valid || err != nil {
		return nil, common.NewError("invalid_signature", "Invalid signature")
	}

	clientID := ctx.Value(constants.ContextKeyClient).(string)
	_ = ctx.Value(constants.ContextKeyClientKey).(string)

	allocationID := allocationObj.ID

	if clientID == "" {
		return nil, common.NewError("invalid_operation", "Invalid client")
	}

	destPath := r.FormValue("dest")
	if destPath == "" {
		return nil, common.NewError("invalid_parameters", "Invalid destination for operation")
	}
	destPath = filepath.Clean(destPath)

	pathHash, err := pathHashFromReq(r, allocationID)
	if err != nil {
		return nil, err
	}

	if clientID == "" || allocationObj.OwnerID != clientID {
		return nil, common.NewError("invalid_operation", "Operation needs to be performed by the owner of the allocation")
	}

	connectionID := r.FormValue("connection_id")
	if connectionID == "" {
		return nil, common.NewError("invalid_parameters", "Invalid connection id passed")
	}

	connectionObj, err := allocation.GetAllocationChanges(ctx, connectionID, allocationID, clientID)
	if err != nil {
		return nil, common.NewError("meta_error", "Error reading metadata for connection")
	}

	objectRef, err := reference.GetLimitedRefFieldsByLookupHash(
		ctx, allocationID, pathHash, []string{"id", "name", "path", "size", "type"})

	if err != nil {
		return nil, common.NewError("invalid_parameters", "Invalid file path. "+err.Error())
	}

	if objectRef.ParentPath == destPath {
		return nil, common.NewError("invalid_parameters", "Invalid destination path. Cannot move to the same parent directory.")
	}
	if objectRef.Type == reference.DIRECTORY {
		isEmpty, err := reference.IsDirectoryEmpty(ctx, objectRef.ID)
		if err != nil {
			return nil, common.NewError("invalid_operation", "Error checking if directory is empty")
		}
		if !isEmpty {
			return nil, common.NewError("invalid_operation", "Directory is not empty")
		}
	}
	newPath := filepath.Join(destPath, objectRef.Name)
	newPath = filepath.Clean(newPath)
	paths, err := common.GetParentPaths(newPath)
	if err != nil {
		return nil, err
	}

	paths = append(paths, newPath)

	refs, err := reference.GetRefsTypeFromPaths(ctx, allocationID, paths)
	if err != nil {
		Logger.Error("Database error", zap.Error(err))
		return nil, common.NewError("database_error", fmt.Sprintf("Got db error while getting refs for %v", paths))
	}

	for _, ref := range refs {
		switch ref.Path {
		case newPath:
			return nil, common.NewError("invalid_parameters", "Invalid destination path. Object Already exists.")
		default:
			if ref.Type == reference.FILE {
				return nil, common.NewError("invalid_path", fmt.Sprintf("%v is of file type", ref.Path))
			}
		}
	}

	allocationChange := &allocation.AllocationChange{}
	allocationChange.ConnectionID = connectionObj.ID
	allocationChange.Size = 0
	allocationChange.LookupHash = pathHash
	allocationChange.Operation = constants.FileOperationMove
	dfc := &allocation.MoveFileChange{
		ConnectionID: connectionObj.ID,
		AllocationID: connectionObj.AllocationID,
		SrcPath:      objectRef.Path,
		DestPath:     newPath,
		Type:         objectRef.Type,
	}
	dfc.SrcPath = objectRef.Path
	connectionObj.AddChange(allocationChange, dfc)

	err = connectionObj.Save(ctx)
	if err != nil {
		Logger.Error("Error in writing the connection meta data", zap.Error(err))
		return nil, common.NewError("connection_write_error", "Error writing the connection meta data")
	}

	result := &allocation.UploadResult{}
	result.Filename = objectRef.Name
	result.Size = objectRef.Size
	return result, nil
}

func (fsh *StorageHandler) DeleteFile(ctx context.Context, r *http.Request, connectionObj *allocation.AllocationChangeCollector) (*allocation.UploadResult, error) {

	path := r.FormValue("path")
	if path == "" {
		return nil, common.NewError("invalid_parameters", "Invalid path")
	}
	fileRef, err := reference.GetLimitedRefFieldsByPath(ctx, connectionObj.AllocationID, path,
		[]string{"path", "name", "size"})

	if err != nil {
		Logger.Error("invalid_file", zap.Error(err))
	}
	_ = ctx.Value(constants.ContextKeyClientKey).(string)
	if fileRef != nil {
		deleteSize := fileRef.Size

		allocationChange := &allocation.AllocationChange{}
		allocationChange.ConnectionID = connectionObj.ID
		allocationChange.Size = 0 - deleteSize
		allocationChange.Operation = constants.FileOperationDelete
		dfc := &allocation.DeleteFileChange{ConnectionID: connectionObj.ID,
			AllocationID: connectionObj.AllocationID, Name: fileRef.Name,
			LookupHash: fileRef.LookupHash, Path: fileRef.Path, Size: deleteSize}

		allocation.UpdateConnectionObjSize(connectionObj.ID, allocationChange.Size)

		connectionObj.AddChange(allocationChange, dfc)

		result := &allocation.UploadResult{}
		result.Filename = fileRef.Name
		result.Hash = fileRef.LookupHash
		result.Size = fileRef.Size

		return result, nil
	}

	return nil, common.NewError("invalid_file", "File does not exist at path")
}

func (fsh *StorageHandler) CreateDir(ctx context.Context, r *http.Request) (*allocation.UploadResult, error) {
	allocationId := ctx.Value(constants.ContextKeyAllocationID).(string)
	allocationTx := ctx.Value(constants.ContextKeyAllocation).(string)
	clientID := ctx.Value(constants.ContextKeyClient).(string)
	allocationObj, err := fsh.verifyAllocation(ctx, allocationId, allocationTx, false)
	if err != nil {
		return nil, common.NewError("invalid_parameters", "Invalid allocation id passed."+err.Error())
	}

	valid, err := verifySignatureFromRequest(allocationTx, r.Header.Get(common.ClientSignatureHeader), r.Header.Get(common.ClientSignatureHeaderV2), allocationObj.OwnerPublicKey)
	if !valid || err != nil {
		return nil, common.NewError("invalid_signature", "Invalid signature")
	}

	allocationID := allocationObj.ID

	if clientID == "" {
		return nil, common.NewError("invalid_operation", "Operation needs to be performed by the owner or the payer of the allocation")
	}

	dirPath := r.FormValue("dir_path")
	if dirPath == "" {
		return nil, common.NewError("invalid_parameters", "Invalid dir path passed")
	}

	customMeta := r.FormValue("custom_meta")

	exisitingRef, err := fsh.checkIfFileAlreadyExists(ctx, allocationID, dirPath)
	if err != nil {
		Logger.Error("Error file reference", zap.Error(err))
	}

	result := &allocation.UploadResult{
		Filename: dirPath,
	}

	if exisitingRef != nil {
		// target directory exists, return StatusOK
		if exisitingRef.Type == reference.DIRECTORY {

			if exisitingRef.CustomMeta != customMeta {
				_ = datastore.GetStore().WithNewTransaction(func(ctx context.Context) error {
					err := reference.UpdateCustomMeta(ctx, exisitingRef, customMeta)
					if err != nil {
						logging.Logger.Error("Error updating custom meta", zap.Error(err))
					}
					return err
				})
			}
			return nil, common.NewError("directory_exists", "Directory already exists`")
		}

		msg := fmt.Sprintf("File at path :%s: already exists", exisitingRef.Path)
		return nil, common.NewError("duplicate_file", msg)
	}
	if !filepath.IsAbs(dirPath) {
		return nil, common.NewError("invalid_path", fmt.Sprintf("%v is not absolute path", dirPath))
	}

	if clientID != allocationObj.OwnerID {
		return nil, common.NewError("invalid_operation", "Operation needs to be performed by the owner or the payer of the allocation")
	}

	if err := validateParentPathType(ctx, allocationID, dirPath); err != nil {
		return nil, err
	}

	connectionID := r.FormValue("connection_id")
	if connectionID == "" {
		return nil, common.NewError("invalid_parameters", "Invalid connection id passed")
	}

	connectionObj, err := allocation.GetAllocationChanges(ctx, connectionID, allocationID, clientID)
	if err != nil {
		return nil, common.NewError("meta_error", "Error reading metadata for connection")
	}

	allocationChange := &allocation.AllocationChange{}
	allocationChange.ConnectionID = connectionObj.ID
	allocationChange.Size = 0
	allocationChange.Operation = constants.FileOperationCreateDir
	var newDir allocation.NewDir
	newDir.ConnectionID = connectionID
	newDir.Path = dirPath
	newDir.AllocationID = allocationID
	newDir.CustomMeta = customMeta

	connectionObj.AddChange(allocationChange, &newDir)

	err = connectionObj.Save(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// WriteFile stores the file into the blobber files system from the HTTP request
func (fsh *StorageHandler) WriteFile(ctx context.Context, r *http.Request) (*allocation.UploadResult, error) {
	startTime := time.Now()
	if r.Method == "GET" {
		return nil, common.NewError("invalid_method", "Invalid method used for the upload URL. Use multi-part form POST / PUT / DELETE / PATCH instead")
	}

	allocationID := ctx.Value(constants.ContextKeyAllocationID).(string)
	allocationTx := ctx.Value(constants.ContextKeyAllocation).(string)
	clientID := ctx.Value(constants.ContextKeyClient).(string)
	logging.Logger.Info("writeFile", zap.String("allocation_id", allocationID))
	connectionID, ok := common.GetField(r, "connection_id")
	if !ok {
		logging.Logger.Error("no_connection_id", zap.String("alloc_id", allocationID))
		return nil, common.NewError("invalid_parameters", "Invalid connection id passed")
	}
	if clientID == "" {
		return nil, common.NewError("invalid_operation", "Operation needs to be performed by the owner or the payer of the allocation")
	}
	if ok := CheckBlacklist(clientID); ok {
		return nil, common.NewError("blacklisted_client", "Client is blacklisted: "+clientID)
	}
	elapsedParseForm := time.Since(startTime)
	st := time.Now()
	if allocation.GetConnectionProcessor(connectionID) == nil {
		allocation.CreateConnectionProcessor(connectionID, allocationID, clientID)
	}

	elapsedGetConnectionProcessor := time.Since(st)
	st = time.Now()

	allocationObj, err := fsh.verifyAllocation(ctx, allocationID, allocationTx, false)
	if err != nil {
		return nil, common.NewError("invalid_parameters", "Invalid allocation id passed."+err.Error())
	}
	logging.Logger.Info("writeFileAllocation")
	if allocationObj.OwnerID != clientID && allocationObj.RepairerID != clientID {
		return nil, common.NewError("invalid_operation", "Operation needs to be performed by the owner or the payer of the allocation")
	}

	elapsedAllocation := time.Since(st)

	if r.Method == http.MethodPost && !allocationObj.CanUpload() {
		return nil, common.NewError("prohibited_allocation_file_options", "Cannot upload data to this allocation.")
	}

	if r.Method == http.MethodPut && !allocationObj.CanUpdate() {
		return nil, common.NewError("prohibited_allocation_file_options", "Cannot update data in this allocation.")
	}

	if r.Method == http.MethodDelete && !allocationObj.CanDelete() {
		return nil, common.NewError("prohibited_allocation_file_options", "Cannot delete data in this allocation.")
	}

	st = time.Now()
	publicKey := allocationObj.OwnerPublicKey

	valid, err := verifySignatureFromRequest(allocationTx, r.Header.Get(common.ClientSignatureHeader), r.Header.Get(common.ClientSignatureHeaderV2), publicKey)

	if !valid || err != nil {
		return nil, common.NewError("invalid_signature", "Invalid signature")
	}

	elapsedVerifySig := time.Since(st)
	st = time.Now()

	cmd := createFileCommand(r)
	err = cmd.IsValidated(ctx, r, allocationObj, clientID)
	if err != nil {
		return nil, err
	}
	elapsedIsValidated := time.Since(st)
	st = time.Now()
	// call process content, which writes to file checks if conn obj needs to be updated and if commit hasher needs to be called
	res, err := cmd.ProcessContent(ctx, allocationObj)
	if err != nil {
		return nil, err
	}
	// Update/Save the change
	if res.UpdateChange {
		_, err := allocation.GetConnectionObj(ctx, connectionID, allocationID, clientID)
		if err != nil {
			return nil, err
		}
		err = cmd.UpdateChange(ctx)
		if err != nil {
			return nil, err
		}
	}

	blocks := cmd.GetNumBlocks()
	if blocks > 0 {
		go AddUploadedData(clientID, blocks)
	}
	elapsedProcessContent := time.Since(st)
	Logger.Info("[upload]elapsed",
		zap.String("alloc_id", allocationID),
		zap.String("file", cmd.GetPath()),
		zap.Duration("parse_form", elapsedParseForm),
		zap.Duration("get_processor", elapsedGetConnectionProcessor),
		zap.Duration("get_alloc", elapsedAllocation),
		zap.Duration("sig", elapsedVerifySig),
		zap.Duration("validate", elapsedIsValidated),
		zap.Duration("process_content", elapsedProcessContent),
		zap.Duration("total", time.Since(startTime)))
	return &res, nil
}

func sanitizeString(input string) string {
	sanitized := strings.ReplaceAll(input, "\n", "")
	sanitized = strings.ReplaceAll(sanitized, "\r", "")
	return sanitized
}

func (fsh *StorageHandler) Rollback(ctx context.Context, r *http.Request) (*blobberhttp.CommitResult, error) {

	startTime := time.Now()
	if r.Method == "GET" {
		return nil, common.NewError("invalid_method", "Invalid method used for the rolllback URL. Use POST instead")
	}

	Logger.Info("Rollback request received")

	allocationId := ctx.Value(constants.ContextKeyAllocationID).(string)
	allocationTx := ctx.Value(constants.ContextKeyAllocation).(string)
	clientID := ctx.Value(constants.ContextKeyClient).(string)
	clientKey := ctx.Value(constants.ContextKeyClientKey).(string)
	clientKeyBytes, _ := hex.DecodeString(clientKey)
	var (
		allocationObj *allocation.Allocation
		err           error
	)

	allocationObj, err = fsh.verifyAllocation(ctx, allocationId, allocationTx, false)
	if err != nil {
		Logger.Error("Error in verifying allocation", zap.Error(err))
		return nil, common.NewError("invalid_parameters", "Invalid allocation id passed."+err.Error())
	}

	if allocationObj.AllocationVersion == 0 {
		Logger.Error("Allocation version is 0", zap.String("allocation_id", allocationObj.ID))
		return nil, common.NewError("invalid_parameters", "Allocation version is not set")
	}

	elapsedAllocation := time.Since(startTime)

	allocationID := allocationObj.ID
	elapsedGetLock := time.Since(startTime) - elapsedAllocation

	if clientID == "" || clientKey == "" {
		return nil, common.NewError("invalid_params", "Please provide clientID and clientKey")
	}

	if allocationObj.OwnerID != clientID || encryption.Hash(clientKeyBytes) != clientID {
		return nil, common.NewError("invalid_operation", "Operation needs to be performed by the owner of the allocation")
	}

	if !allocationObj.IsRedeemRequired {
		return nil, common.NewError("invalid_operation", "Last commit is rollback")
	}

	versionMarkerString := r.FormValue("version_marker")
	versionMarker := writemarker.VersionMarker{}
	err = json.Unmarshal([]byte(versionMarkerString), &versionMarker)
	if err != nil {
		return nil, common.NewErrorf("invalid_parameters",
			"Invalid parameters. Error parsing the writemarker for commit: %v",
			err)
	}

	if versionMarker.IsRepair {
		return nil, common.NewError("invalid_parameters", "Invalid version marker passed. Rollback marker cannot be a repair marker")
	}

	var result blobberhttp.CommitResult
	err = versionMarker.Verify(allocationID, allocationObj.OwnerPublicKey)
	if err != nil {
		return nil, common.NewError("invalid_parameters", "Invalid version marker passed: "+err.Error())
	}

	if versionMarker.Version == allocationObj.AllocationVersion {
		return nil, common.NewError("invalid_parameters", "Invalid version marker passed. Version marker is same as the current version")
	}

	currentVersionMarker, err := writemarker.GetCurrentVersion(ctx, allocationID)
	if err != nil {
		return nil, common.NewError("invalid_parameters", "Error getting the current version marker")
	}
	if currentVersionMarker.IsRepair {
		return nil, common.NewError("invalid_parameters", "Invalid version marker passed. Allocation is in repair mode")
	}
	if versionMarker.Version != currentVersionMarker.Version-1 {
		return nil, common.NewError("invalid_parameters", "Invalid version marker passed. Version marker is not the previous version")
	}

	elapsedWritePreRedeem := time.Since(startTime) - elapsedAllocation - elapsedGetLock
	timeoutCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()
	c := datastore.GetStore().CreateTransaction(timeoutCtx)
	txn := datastore.GetStore().GetTransaction(c)
	err = allocation.ApplyRollback(c, allocationID, allocationObj.AllocationVersion)
	if err != nil {
		txn.Rollback()
		return nil, common.NewError("allocation_rollback_error", "Error applying the rollback for allocation: "+err.Error())
	}
	elapsedApplyRollback := time.Since(startTime) - elapsedAllocation - elapsedGetLock - elapsedWritePreRedeem

	alloc, err := allocation.Repo.GetByIdAndLock(c, allocationID)
	Logger.Info("[rollback]Lock Allocation", zap.Bool("is_redeem_required", alloc.IsRedeemRequired), zap.String("allocation_root", alloc.AllocationRoot), zap.String("latest_wm_redeemed", alloc.LatestRedeemedWM))
	if err != nil {
		txn.Rollback()
		return &result, common.NewError("allocation_read_error", "Error reading the allocation object")
	}

	alloc.BlobberSizeUsed = alloc.PrevBlobberSizeUsed
	alloc.UsedSize = alloc.PrevUsedSize
	alloc.IsRedeemRequired = false
	alloc.AllocationVersion = versionMarker.Version
	updateMap := map[string]interface{}{
		"blobber_size_used":  alloc.BlobberSizeUsed,
		"used_size":          alloc.UsedSize,
		"is_redeem_required": false,
		"allocation_version": versionMarker.Version,
	}

	updateOption := func(a *allocation.Allocation) {
		a.BlobberSizeUsed = alloc.BlobberSizeUsed
		a.UsedSize = alloc.UsedSize
		a.AllocationVersion = alloc.AllocationVersion
		a.FileMetaRoot = alloc.FileMetaRoot
		a.IsRedeemRequired = alloc.IsRedeemRequired
	}
	err = txn.Create(&versionMarker).Error
	if err != nil {
		txn.Rollback()
		return &result, common.NewError("write_marker_error", "Error persisting the write marker "+err.Error())
	}
	if err = allocation.Repo.UpdateAllocation(c, alloc, updateMap, updateOption); err != nil {
		txn.Rollback()
		return &result, common.NewError("allocation_write_error", "Error persisting the allocation object "+err.Error())
	}

	err = txn.Commit().Error
	if err != nil {
		return &result, common.NewError("allocation_commit_error", "Error committing the transaction "+err.Error())
	}
	err = allocation.CommitRollback(allocationID)
	if err != nil {
		Logger.Error("Error committing the rollback for allocation", zap.Error(err))
	}

	elapsedCommitRollback := time.Since(startTime) - elapsedAllocation - elapsedGetLock - elapsedWritePreRedeem
	result.Success = true
	result.ErrorMessage = ""
	commitOperation := "rollback"

	Logger.Info("[rollback]"+commitOperation,
		zap.String("alloc_id", allocationID),
		zap.Duration("get_alloc", elapsedAllocation),
		zap.Duration("get-lock", elapsedGetLock),
		zap.Duration("write-pre-redeem", elapsedWritePreRedeem),
		zap.Duration("apply-rollback", elapsedApplyRollback),
		zap.Duration("total", time.Since(startTime)),
		zap.Duration("commit-rollback", elapsedCommitRollback),
	)

	return &result, nil
}
