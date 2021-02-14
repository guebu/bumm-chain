package state

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/guebu/common-utils/logger"
	"go.mod/helper"
	"go.mod/model"
	"go.mod/model/account"
	"go.mod/model/trx"
	"os"
)

type Hash [32]byte

type State struct {
	Balances 	map[account.Account]uint
	txMempool 	[]trx.Trx
	dbFile 		*os.File
	latestBlockHash 	Hash
}

func (s *State) GetLatestBlockHash() Hash {
	return s.latestBlockHash
}

func (s *State) apply(tx trx.Trx) error {

	logger.Info("Appliying trx to state", "Layer:Model", "Func:apply", "Status:Start")
	if tx.IsReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}

	if tx.Value > s.Balances[tx.From] {
		return fmt.Errorf("insufficient balance")
	}

	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value

	logger.Info("End of appliying trx to state", "Layer:Model", "Func:apply", "Status:End")
	return nil
}

/*
func (s *State) GetSnapshot() (*Hash, error) {
	if err := s.doSnapshot(); err != nil {
		return nil, err
	}
	return &s.latestBlockHash, nil
}
 */

func (s *State) Add(trx trx.Trx) error {
	logger.Info("Start adding trx to mem pool!", "Layer:Model", "Func:Add", "Status:Start")
	if err := s.apply(trx); err != nil {
		logger.Error("Error in applying trx to current state!", err, "Layer:Model", "Func:NewStateFromDisk", "Status:Error")
		return err
	}
	s.txMempool = append(s.txMempool, trx)
	logger.Info("Added trx to mem pool successfully!", "Layer:Model", "Func:Add", "Status:End")
	fmt.Println(s.txMempool)
	return nil
}

func (s *State) Persist() (*Hash, error) {
	logger.Info("Start of persisting mempool!", "Layer:Model", "Func:Persist", "Status:Start")

	// Create a new Block with ONLY the new TXs
	block := NewBlock(
		s.latestBlockHash,
		//int64(time.Now().Unix()),
		int64(0),
		s.txMempool,
	)

	// Compute the hash of the new block
	blockHash, err := block.Hash()
	logger.Info(fmt.Sprintf("Generated hash for new block: %x", blockHash), "Layer:Model", "Func:Persist", "Status:Pending")
	if err != nil {
		logger.Error("Hashing of block was not successfull!", err, "Layer:Model", "Func:Persist", "Status:Error")
		return nil, err
	}

	blockFs := BlockFS{	 *blockHash,
						*block }

	// Encode it into a JSON string
	blockFsJson, err := json.Marshal(blockFs)
	if err != nil {
		logger.Error("Marshaling of block was not successfull!", err, "Layer:Model", "Func:Persist", "Status:Error")
		return nil, err
	}

	logger.Info("###############################", "Layer:Model", "Func:Persist", "Status:Pending")
	logger.Info("Persisting new block to disk!", "Layer:Model", "Func:Persist", "Status:Pending")
	logger.Info(fmt.Sprintf("\t%s\n", blockFsJson), "Layer:Model", "Func:Persist", "Status:Pending")
	logger.Info(fmt.Sprintf("Hash of block: %x", blockHash), "Layer:Model", "Func:Persist", "Status:Pending")
	logger.Info(fmt.Sprintf("Hash of parent: %x", block.Header.Parent), "Layer:Model", "Func:Persist", "Status:Pending")
	logger.Info("###############################", "Layer:Model", "Func:Persist", "Status:Pending")

	// Write it to the DB file on a new line
	if _, err = s.dbFile.Write(append(blockFsJson, '\n')); err != nil {
		logger.Error("Error during writing trx data to file!", err, "Layer:Model", "Func:Persist", "Status:Error")
		return nil, err
	}
	s.latestBlockHash = *blockHash

	// Reset the Mempool
	s.txMempool = []trx.Trx{}

	logger.Info("Persisted mempool to disk!", "Layer:Model", "Func:Persist", "Status:End")

	return &s.latestBlockHash, nil


}

func NewStateFromDisk() (*State, error) {
	logger.Info("Start reading state from disk!", "Layer:Model", "Func:NewStateFromDisk", "Status:Start")

	singularityFilePath := helper.GetSingularityFilePath()

	gen, err := model.LoadSingularity(singularityFilePath)

	if err != nil {
		logger.Error("Error in loading singularity file", err, "Layer:Model", "Func:NewStateFromDisk", "Status:Error")
		return nil, err
	}
	logger.Info("Singularity File read successfull", "Layer:Model", "Func:NewStateFromDisk", "Status:Pending")
	balances := make(map[account.Account]uint)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	txDbFilePath := helper.GetTrxDBFilePath()
	dbFile, err := os.OpenFile(txDbFilePath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		logger.Error("Error in opening Trx-DB file", err, "Layer:Model", "Func:NewStateFromDisk", "Status:Error")
		return nil, err
	}

	scanner := bufio.NewScanner(dbFile)

	state := &State{balances, make([]trx.Trx, 0), dbFile, Hash{}}
	noOfEntries := 0
	// Iterate over each the tx.db file's line
	for scanner.Scan() {
		logger.Info(fmt.Sprintf("Number of entries: %v", noOfEntries), "Layer:Model", "Func:NewStateFromDisk", "Status:Pending")

		if err := scanner.Err(); err != nil {
			logger.Error("Error while scanning opened Trx-DB file", err, "Layer:Model", "Func:NewStateFromDisk", "Status:Error")
			return nil, err
		}
		// Convert JSON encoded TX into an object (struct)
		var blockFs BlockFS
		blockFSJson := scanner.Bytes()
		logger.Info("Scanned bytes: --------------------------------", "Layer:Model", "Func:NewStateFromDisk", "Status:Pending")
		logger.Info(string(blockFSJson),"Layer:Model", "Func:NewStateFromDisk", "Status:Pending")
		logger.Info("Scanned bytes: --------------------------------", "Layer:Model", "Func:NewStateFromDisk", "Status:Pending")
		if err := json.Unmarshal(blockFSJson, &blockFs); err != nil {
			logger.Error("Error while unmarshaling trx information read from DB-file!", err, "Layer:Model", "Func:NewStateFromDisk", "Status:Error")
			return nil, err
		}
		logger.Info(blockFs.ToString(), "Layer:Model", "Func:NewStateFromDisk", "Status:Pending")
		logger.Info(fmt.Sprintf("Block-Hash: %x", blockFs.Key) , "Layer:Model", "Func:NewStateFromDisk", "Status:Pending")


		trxs := blockFs.Value.TRXs

		for i, trx := range trxs {
			// Rebuild the state (user balances),
			// as a series of events
			if err := state.apply(trx); err != nil {
				logger.Error("Error while applying trx to existing state...", err, "Layer:Model", "Func:NewStateFromDisk", "Status:Error")
				return nil, err
			}

			logger.Info(fmt.Sprintf("Successfully processed transaction no. %d", i))
		}
		noOfEntries++

		state.latestBlockHash = blockFs.Key
	}

	logger.Info("Finished reading state from disk!", "Layer:Model", "Func:NewStateFromDisk", "Status:End")
	return state, nil
}

func (s *State) Close() error {
	return s.dbFile.Close()
}

/*
func (s *State) doSnapshot() error {
	logger.Info("Start creating a hash/snapshot for current state!", "Layer:Model", "Func:doSnapshot", "Status:Start")
	// Re-read the whole file from the first byte
	_, err := s.dbFile.Seek(0, 0)
	if err != nil {
		logger.Error("Error during seeking DB-File!", err, "Layer:Model", "Func:doSnapshot", "Status:Error")
		return err
	}
	txsData, err := ioutil.ReadAll(s.dbFile)
	if err != nil {
		logger.Error("Error in reading info from DB-File!", err, "Layer:Model", "Func:doSnapshot", "Status:Error")
		return err
	}
	s.latestBlockHash = sha256.Sum256(txsData)
	logger.Info("Snapshot/Hash successfully created for current state!", "Layer:Model", "Func:doSnapshot", "Status:End")
	return nil
}
 */


func (s *State) AddBlock(b Block) error {
	logger.Info("Start adding Block to current state/blockchain!", "Layer:Model", "Func:AddBlock", "Status:Start")
	for _, trx := range b.TRXs {
		if err := s.Add(trx); err != nil {
			logger.Error("Error in adding trx to state", err, "Layer:Model", "Func:AddBlock", "Status:Error")
			return err
		}
	}
	logger.Info("Block added successfully to current state/blockchain!", "Layer:Model", "Func:AddBlock", "Status:End")
	return nil
}
