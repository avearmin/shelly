package cmd

import (
	"fmt"
	"os"

	"github.com/avearmin/shelly/internal/cmdstore"
	"github.com/avearmin/shelly/internal/configstore"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(delCmd)
}

var delCmd = &cobra.Command{
	Use:   "del",
	Short: "Delete an alias and shell command.",
	Long:  "Delete an alias and shell command that is managed by shelly.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configstore.MustHaveConfig()

		config, err := configstore.Load()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		cmds, err := cmdstore.Load(config.CmdsPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		delete(cmds, args[0])

		if err := cmdstore.Save(config.CmdsPath, cmds); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		os.Exit(0)
	},
}
