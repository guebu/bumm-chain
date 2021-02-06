package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.mod/config"
)

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Describes version.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s.%s.%s-beta %s", config.Major, config.Minor, config.Fix, config.Verbal)
		},
	}

