package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var flagYaml,
	flagYml bool

var configCmd = &cobra.Command{
	Use: "config",
	Aliases: []string{
		"c",
	},
	Short: "Get the current running config",
	Long:  `Get the current running config`,
	Run: func(cmd *cobra.Command, args []string) {
		var settings []byte
		if flagYaml || flagYml {
			settings, _ = yaml.Marshal(viper.AllSettings())
		} else {
			settings, _ = json.MarshalIndent(viper.AllSettings(), "", "  ")
		}
		fmt.Println(string(settings))
	},
}

func init() {
	configCmd.Flags().BoolVarP(&flagYaml, "yaml", "y", false, "output yaml instead of json")
	configCmd.Flags().BoolVar(&flagYml, "yml", false, "alias for -y and --yaml")
	RootCmd.AddCommand(configCmd)
}
