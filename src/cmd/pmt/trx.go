package main

import (
	"fmt"
	"github.com/guebu/common-utils/errors"
	"github.com/spf13/cobra"
	"go.mod/config"
	"go.mod/model/account"
	"go.mod/model/state"
	"go.mod/model/trx"
	"os"
)

func trxCmd() *cobra.Command {
	var txsCmd = &cobra.Command{
		Use: "trx",
		Short: "Interact with trxs (add...).",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return errors.NewBadRequestError("Uncorrect Usage!", nil)
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	txsCmd.AddCommand(trxAddCmd())
	return txsCmd
}

func trxAddCmd() *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "add",
		Short: "Adds new TX to database.",
		Run: func(cmd *cobra.Command, args []string) {
			from, _ := cmd.Flags().GetString(config.FromCmdKey)
			to, _ := cmd.Flags().GetString(config.ToCmdKey)
			value, _ := cmd.Flags().GetUint(config.ValueCmdKey)

			fromAcc := account.NewAccount(from)
			toAcc := account.NewAccount(to)

			trx := trx.NewTrx(*fromAcc, *toAcc, value, "")

			state, err := state.NewStateFromDisk()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			// defer means, at the end of this function execution,
			// execute the following statement (close DB file with all TXs)
			defer state.Close()
			// Add the TX to an in-memory array (pool)
			err = state.Add(*trx)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			// Flush the mempool TXs to disk
			err = state.Persist()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Println("TX successfully added to the ledger.")
		},
	}

	cmd.Flags().String(config.FromCmdKey, "", "From what account to send tokens")
	cmd.MarkFlagRequired(config.FromCmdKey)
	cmd.Flags().String(config.ToCmdKey, "", "To what account to send tokens")
	cmd.MarkFlagRequired(config.ToCmdKey)
	cmd.Flags().Uint(config.ValueCmdKey, 0, "How many tokens to send")
	cmd.MarkFlagRequired(config.ValueCmdKey)
	return cmd
}