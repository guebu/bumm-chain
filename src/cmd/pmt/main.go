
package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	var tbbCmd = &cobra.Command{
		Use: "pmt",
		Short: "The Blockchain Bar CLI",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// add command to
	tbbCmd.AddCommand(versionCmd)
	tbbCmd.AddCommand(balancesListCmd)
	tbbCmd.AddCommand(trxCmd())
	tbbCmd.AddCommand(migCmd)

	fmt.Println("----------------------------------")
	fmt.Println("Start main routine of commands...")
	fmt.Println("----------------------------------")
	err := tbbCmd.Execute()
	if err != nil {
		fmt.Println("----------------------------------------")
		fmt.Println("Error while trying to execute command...")
		fmt.Println("----------------------------------------")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
