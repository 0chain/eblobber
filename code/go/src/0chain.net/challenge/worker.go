package challenge

import (
	"container/list"
	"context"
	"encoding/json"
	"sync"
	"time"

	"0chain.net/chain"
	"0chain.net/config"
	"0chain.net/datastore"
	"0chain.net/filestore"
	"0chain.net/lock"
	. "0chain.net/logging"
	"0chain.net/node"
	"0chain.net/stats"
	"0chain.net/transaction"
	"0chain.net/writemarker"

	"go.uber.org/zap"
)

type BCChallengeResponse struct {
	BlobberID  string             `json:"blobber_id"`
	Challenges []*ChallengeEntity `json:"challenges"`
}

var dataStore datastore.Store
var fileStore filestore.FileStore

func SetupWorkers(ctx context.Context, metaStore datastore.Store, fsStore filestore.FileStore) {
	dataStore = metaStore
	fileStore = fsStore
	go FindChallenges(ctx)
}

func RespondToChallenge(challengeID string) {
	ctx := context.Background()
	newctx := dataStore.WithConnection(ctx)
	challengeObj := Provider().(*ChallengeEntity)
	challengeObj.ID = challengeID
	err := challengeObj.Read(newctx, challengeObj.GetKey())
	if err != nil {
		Logger.Error("Error reading challenge from the database.", zap.Error(err), zap.String("challenge_id", challengeID))
	}

	mutex := lock.GetMutex(challengeObj.GetKey())
	mutex.Lock()
	if challengeObj.Status == Error {
		challengeObj.Retries++
	}
	err = challengeObj.SendDataBlockToValidators(newctx, fileStore)
	if err != nil {
		Logger.Error("Error in responding to challenge. ", zap.Any("error", err.Error()))
	}

	err = dataStore.Commit(newctx)

	if err != nil {
		Logger.Error("Error in challenge commit to DB", zap.Error(err), zap.String("challenge_id", challengeID))
	}

	mutex.Unlock()

	if challengeObj.ObjectPath != nil && challengeObj.Status == Committed && challengeObj.ObjectPath.FileBlockNum > 0 {
		stats.FileChallenged(newctx, challengeObj.AllocationID, challengeObj.ObjectPath.Meta["path"].(string), challengeObj.CommitTxnID)
	}
	newctx.Done()
	challengeWorker.Done()
	Logger.Info("Challenge has been processed", zap.Any("id", challengeObj.ID), zap.String("txn", challengeObj.CommitTxnID))
}

var challengeHandler = func(ctx context.Context, key datastore.Key, value []byte) error {
	challengeObj := Provider().(*ChallengeEntity)
	err := json.Unmarshal(value, challengeObj)
	if err != nil {
		return err
	}

	if challengeObj.Status != Committed && challengeObj.Status != Failed && challengeObj.Retries < 20 {
		unredeemedMarkers.PushBack(challengeObj.ID)
	}
	return nil
}

var challengeWorker sync.WaitGroup
var numOfWorkers = 0
var iterInprogress = false
var unredeemedMarkers *list.List

func FindChallenges(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(config.Configuration.ChallengeResolveFreq) * time.Second)
	for true {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if !iterInprogress && numOfWorkers == 0 {
				unredeemedMarkers = list.New()
				iterInprogress = true
				rctx := dataStore.WithReadOnlyConnection(context.Background())
				dataStore.IteratePrefix(rctx, "challenge:", challengeHandler)
				dataStore.Discard(rctx)
				rctx.Done()
				for e := unredeemedMarkers.Front(); e != nil; e = e.Next() {
					if numOfWorkers < config.Configuration.ChallengeResolveNumWorkers {
						numOfWorkers++
						challengeWorker.Add(1)
						Logger.Info("Starting challenge with ID: " + e.Value.(string))
						go RespondToChallenge(e.Value.(string))
					} else {
						challengeWorker.Wait()
					}
				}
				if numOfWorkers > 0 {
					challengeWorker.Wait()
				}

				iterInprogress = false
				numOfWorkers = 0
				params := make(map[string]string)
				params["blobber"] = node.Self.ID
				var blobberChallenges BCChallengeResponse
				blobberChallenges.Challenges = make([]*ChallengeEntity, 0)
				_, err := transaction.MakeSCRestAPICall(transaction.STORAGE_CONTRACT_ADDRESS, "/openchallenges", params, chain.GetServerChain(), &blobberChallenges)
				if err == nil {
					tCtx := dataStore.WithConnection(ctx)
					for _, v := range blobberChallenges.Challenges {
						if v == nil {
							Logger.Info("No challenge entity from the challenge map")
							continue
						}
						challengeObj := v
						err = challengeObj.Read(tCtx, v.GetKey())
						if err == datastore.ErrKeyNotFound {
							Logger.Info("Adding new challenge found from blockchain", zap.String("challenge", v.ID))
							writeMarkerEntity := writemarker.Provider().(*writemarker.WriteMarkerEntity)
							writeMarkerEntity.WM = &writemarker.WriteMarker{AllocationID: challengeObj.AllocationID, AllocationRoot: challengeObj.AllocationRoot}

							err = writeMarkerEntity.Read(tCtx, writeMarkerEntity.GetKey())
							if err != nil {
								continue
							}
							challengeObj.WriteMarker = writeMarkerEntity.GetKey()
							challengeObj.ValidationTickets = make([]*ValidationTicket, len(challengeObj.Validators))
							challengeObj.Write(tCtx)
						}
					}
					dataStore.Commit(tCtx)
					tCtx.Done()
				}

			}
		}
	}

}