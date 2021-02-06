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

type State struct {
	Balances map[account.Account]uint
	txMempool []trx.Trx
	dbFile *os.File
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

func (s *State) Persist() error {
	logger.Info("Start of persisting mempool!", "Layer:Model", "Func:Persist", "Status:Start")

	// Make a copy of mempool because the s.txMempool will be modified
	// in the loop below
	mempool := make([]trx.Trx, len(s.txMempool))
	copy(mempool, s.txMempool)

	for i := 0; i < len(mempool); i++ {
		txJson, err := json.Marshal(mempool[i])
		if err != nil {
			logger.Error("Error in marshalling mempool content!", err, "Layer:Model", "Func:Persist", "Status:Error")
			return err
		}
		if _, err = s.dbFile.Write(append(txJson, '\n')); err != nil {
			logger.Error("Error in appending trx from mempool to file!", err, "Layer:Model", "Func:Persist", "Status:Error")
			return err
		}
		// Remove the TX written to a file from the mempool
		// Yes... this particular Go syntax is a bit weird
		s.txMempool = append(s.txMempool[:0], s.txMempool[0+1:]...)
	}
	logger.Info("End of persisting mempool!", "Layer:Model", "Func:Persist", "Status:End")
	return nil
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
	state := &State{balances, make([]trx.Trx, 0), dbFile}
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