package service

import (
	"go.mod/model/state"
)

func NewStateFromDisk() (*state.State, error) {
	/*
	// get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	genFilePath := filepath.Join(cwd, "database", "genesis.json")
	gen, err := loadGenesis(genFilePath)
	if err != nil {
		return nil, err
	}
	balances := make(map[model.Account]uint)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	txDbFilePath := filepath.Join(cwd, "database", "tx.db")
	f, err := os.OpenFile(txDbFilePath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	state := &State{balances, make([]Tx, 0), f}
	// Iterate over each the tx.db file's line
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		// Convert JSON encoded TX into an object (struct)
		var tx Tx
		json.Unmarshal(scanner.Bytes(), &tx)
		// Rebuild the state (user balances),
		// as a series of events
		if err := state.apply(tx); err != nil {
			return nil, err
		}
	}
	*)
	 */
	return nil, nil
}
