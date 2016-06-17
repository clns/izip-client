package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "0.1.1"

// VersionCmd prints the current version.
var VersionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Print the version of this tool",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s %v\n", NAME, Version)
	},
}

func init() {
	RootCmd.AddCommand(VersionCmd)
}
