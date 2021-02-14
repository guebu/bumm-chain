package main

import (
	"fmt"
	"github.com/guebu/common-utils/logger"
	"github.com/spf13/cobra"
	"go.mod/model/state"
	"go.mod/model/trx"
	"os"
	"time"

)


var migCmd = &cobra.Command{
	Use: "mig",
	Short: "Migrate the old tx.db in the new format!",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Start Migration Command!", "Layer:Cmd", "Func:migCmd", "Status:Start")

		myState, err := state.NewStateFromDisk()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer myState.Close()

		block0 := *state.NewBlock(
			state.Hash{},
			int64(time.Now().Unix()),
			[]trx.Trx{
				*trx.NewTrx("guebu","usi",3,""),
				*trx.NewTrx("guebu","quirin",7,"pocket money"),
				*trx.NewTrx("guebu","ferdl",2,"pocket money"),
				*trx.NewTrx("guebu","guebu",1,"Reward"),
				*trx.NewTrx("guebu","ferdl",10,"trx1"),
				*trx.NewTrx("guebu","ferdl",10,"trx1"),
				*trx.NewTrx("guebu","ferdl",10,"trx1"),
				*trx.NewTrx("guebu","ferdl",10,"trx1"),
				*trx.NewTrx("guebu","ferdl",10,"trx1"),
				*trx.NewTrx("guebu","ferdl",10,"trx1"),
				*trx.NewTrx("guebu","ferdl",10,"trx1"),
				*trx.NewTrx("guebu","ferdl",10,"trx2"),
				*trx.NewTrx("guebu","ferdl",10,"trx1"),
				*trx.NewTrx("guebu","ferdl",10,"trx2"),
				*trx.NewTrx("guebu","ferdl",10,"trx1"),
				*trx.NewTrx("guebu","ferdl",10,"trx2"),
				*trx.NewTrx("guebu","ferdl",10,"trx1"),
				*trx.NewTrx("guebu","ferdl",10,"trx2"),
				*trx.NewTrx("guebu","ferdl",10,"trx1"),
				*trx.NewTrx("guebu","ferdl",10,"trx2"),
				*trx.NewTrx("guebu","ferdl",10,"trx1"),
				*trx.NewTrx("guebu","ferdl",10,"trx2"),
				*trx.NewTrx("guebu","ferdl",3,""),
				*trx.NewTrx("guebu","quirin",3,""),
			},
		)

		myState.AddBlock(block0)
		block0hash, _ := myState.Persist()

		block1 := *state.NewBlock(
			*block0hash,
			int64(time.Now().Unix()),
			[]trx.Trx{
				*trx.NewTrx("guebu","guebu",10,"Reward"),
				*trx.NewTrx("guebu","usi",10,"pocket money"),
				*trx.NewTrx("guebu","ferdl",10,"pocket money"),
				*trx.NewTrx("guebu","quirin",15,"Reward"),
				*trx.NewTrx("guebu","quirin",15,"trx1"),
				*trx.NewTrx("guebu","quirin",15,"trx1"),
				*trx.NewTrx("guebu","quirin",15,"trx1"),
				*trx.NewTrx("guebu","quirin",15,"trx1"),
				*trx.NewTrx("guebu","quirin",15,"trx1"),
				*trx.NewTrx("guebu","usi",15,"trx1"),
				*trx.NewTrx("guebu","usi",15,"trx1"),
				*trx.NewTrx("guebu","usi",15,"trx2"),
			},
		)

		myState.AddBlock(block1)
		block1hash, _ := myState.Persist()

		logger.Info("Data successfully migrated!", "App:Migration", "Layer:migration.go", "Func:main", "Status:End")
		logger.Info(fmt.Sprintf("%x", block1hash), "App:Migration", "Layer:migration.go", "Func:main", "Status:End")



		logger.Info("Start Migration Command!", "Layer:Cmd", "Func:migCmd", "Status:End")
	},
}


