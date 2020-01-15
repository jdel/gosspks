package cmd // import jdel.org/gosspks/cmd

import (
	"fmt"

	"jdel.org/gosspks/cfg"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use: "version",
	Aliases: []string{
		"v",
	},
	Short: "Get the version of sspks",
	Long:  `Get the version of sspks`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gosspks version: ", cfg.Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
