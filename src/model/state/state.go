package state

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/guebu/common-utils/logger"
	"go.mod/helper"
	"go.mod/model"
	"go.mod/model/account"
	"go.mod/model/trx"
	"io/ioutil"
	"os"
)

type Snapshot [32]byte

type State struct {
	Balances 	map[account.Account]uint
	txMempool 	[]trx.Trx
	dbFile 		*os.File
	snapshot 	Snapshot
}

func (s *State) apply(tx trx.Trx) error {

	if tx.IsReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}

	if tx.Value > s.Balances[tx.From] {
		return fmt.Errorf("insufficient balance")
	}

	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value

	return nil
}

func (s *State) GetSnapshot() (*Snapshot, error) {
	if err := s.doSnapshot(); err != nil {
		return nil, err
	}
	return &s.snapshot, nil
}

func (s *State) Add(trx trx.Trx) error {
	logger.Info("Start addint trx to mem pool!", "Layer:Model", "Func:Add", "Status:Start")
	if err := s.apply(trx); err != nil {
		logger.Error("Error in applying trx to current state!", err, "Layer:Model", "Func:NewStateFromDisk", "Status:Error")
		return err
	}
	s.txMempool = append(s.txMempool, trx)
	logger.Info("Added trx to mem pool successfully!", "Layer:Model", "Func:Add", "Status:End")
	fmt.Println(s.txMempool)
	return nil
}

func (s *State) Persist() (*Snapshot, error) {
	logger.Info("Start of persisting mempool!", "Layer:Model", "Func:Persist", "Status:Start")

	// Make a copy of mempool because the s.txMempool will be modified
	// in the loop below
	mempool := make([]trx.Trx, len(s.txMempool))
	copy(mempool, s.txMempool)

	for i := 0; i < len(mempool); i++ {
		txJson, err := json.Marshal(mempool[i])
		if err != nil {
			logger.Error("Error in marshalling mempool content!", err, "Layer:Model", "Func:Persist", "Status:Error")
			return nil, err
		}

		logger.Info(fmt.Sprintf("Persisting new trx to disk: %s", txJson), "Layer:Model", "Func:Persist", "Status:Pending")

		if _, err = s.dbFile.Write(append(txJson, '\n')); err != nil {
			logger.Error("Error in appending trx from mempool to file!", err, "Layer:Model", "Func:Persist", "Status:Error")
			return nil, err
		}

		// Compute Snapshot
		if err := s.doSnapshot(); err != nil {
			logger.Error("Error while computing hast for DB!", err, "Layer:Model", "Func:Persist", "Status:Error")
			return nil, err
		}

		logger.Info(fmt.Sprintf("New DB Snapshot:  %x", s.snapshot), "Layer:Model", "Func:Persist", "Status:Pending")
		fmt.Printf("New DB Snapshot: %x\n", s.snapshot)
		// Remove the TX written to a file from the mempool
		// Yes... this particular Go syntax is a bit weird
		s.txMempool = append(s.txMempool[:0], s.txMempool[0+1:]...)
	}
	logger.Info("End of persisting mempool!", "Layer:Model", "Func:Persist", "Status:End")
	return &s.snapshot, nil
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

	var byteArray [32]byte
	var initialHash Snapshot = byteArray

	state := &State{balances, make([]trx.Trx, 0), dbFile, initialHash}
	// Iterate over each the tx.db file's line
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			logger.Error("Error while scanning opened Trx-DB file", err, "Layer:Model", "Func:NewStateFromDisk", "Status:Error")
			return nil, err
		}
		// Convert JSON encoded TX into an object (struct)
		var trx trx.Trx
		if err := json.Unmarshal(scanner.Bytes(), &trx); err != nil {
			logger.Error("Error while unmarshaling trx information read from DB-file!", err, "Layer:Model", "Func:NewStateFromDisk", "Status:Error")
			return nil, err
		}

		// Rebuild the state (user balances),
		// as a series of events
		if err := state.apply(trx); err != nil {
			logger.Error("Error while applying trx to existing state...", err, "Layer:Model", "Func:NewStateFromDisk", "Status:Error")
			return nil, err
		}
	}
	logger.Info("Finished reading state from disk!", "Layer:Model", "Func:NewStateFromDisk", "Status:End")
	return state, nil
}

func (s *State) Close() error {
	return s.dbFile.Close()
}

func (s *State) doSnapshot() error {
	// Re-read the whole file from the first byte
	_, err := s.dbFile.Seek(0, 0)
	if err != nil {
		return err
	}
	txsData, err := ioutil.ReadAll(s.dbFile)
	if err != nil {
		return err
	}
	s.snapshot = sha256.Sum256(txsData)
	return nil
}