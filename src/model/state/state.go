package state

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/guebu/common-utils/logger"
	"go.mod/helper"
	"go.mod/model"
	"go.mod/model/trx"
	"os"
)

type State struct {
	Balances map[model.Account]uint
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

func NewStateFromDisk() (*State, error) {
	logger.Info("Start reading state from disk!", "Layer:Model", "Func:NewStateFromDisk", "Status:Start")
	// get current working directory
	/*
	cwd, err := os.Getwd()
	if err != nil {
		logger.Error("Error in getting current working directory", err, "Layer:Model", "Func:NewStateFromDisk", "Status:Error")
		return nil, err
	}
	 */
	singularityFilePath := helper.GetSingularityFilePath()

	gen, err := model.LoadSingularity(singularityFilePath)

	if err != nil {
		logger.Error("Error in loading singularity file", err, "Layer:Model", "Func:NewStateFromDisk", "Status:Error")
		return nil, err
	}
	logger.Info("Singularity File read successfull", "Layer:Model", "Func:NewStateFromDisk", "Status:Pending")
	balances := make(map[model.Account]uint)
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