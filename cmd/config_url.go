package cmd

import (
	"fmt"

	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ConfigURLCmd = &cobra.Command{
	Use:     "url [<value>]",
	Short:   "Print or save the url from/into the config file",
	Example: `  $ ` + NAME + ` config url https://example.com`,
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := NormalizeURL(viper.GetString("url"))
		if len(args) == 0 {
			fmt.Println(u)
			os.Exit(0)
		}
		out, err := NormalizeURL(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			return
		}
		viper.Set("url", out)
		if err := SaveViperConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	ConfigCmd.AddCommand(ConfigURLCmd)
}
