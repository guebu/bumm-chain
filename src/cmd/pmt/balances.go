package main

import (
	"fmt"
	"github.com/guebu/common-utils/errors"
	"github.com/guebu/common-utils/logger"
	"github.com/spf13/cobra"
	"go.mod/model/state"
	"os"
)

func balancesCmd() *cobra.Command {
	var balancesCmd = &cobra.Command{
		Use: "balances",
		Short: "Interact with balances (list...).",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return errors.NewBadRequestError("Incorrect usage!", nil)
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	balancesCmd.AddCommand(balancesListCmd)
	return balancesCmd
}

var balancesListCmd = &cobra.Command{
	Use: "list",
	Short: "Lists all balances.",
	Run: func(cmd *cobra.Command, args []string) {
		state, err := state.NewStateFromDisk()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer state.Close()

		snapshot, err := state.GetSnapshot()

		if err != nil {
			logger.Error("Snapshot couldn't be created!", err, "Layer:Cmd", "Status:Error")
		}
		fmt.Println("__________________")
		fmt.Printf("Snapshot: %x\n", snapshot)
		fmt.Println("__________________")
		fmt.Println("Accounts balances:")
		fmt.Println("__________________")

		for account, balance := range state.Balances {
			fmt.Println(fmt.Sprintf("%s: %d", account, balance))
		}
	},
}