package cmd

import (
	"fmt"
	"github.com/avearmin/shelly/internal/cmdstore"
	"github.com/avearmin/shelly/internal/configstore"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an alias and shell command.",
	Long: `Add an alias and shell command to be managed by shelly. You can
		delete this command using 'shelly del'.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !configstore.Exists() {
			fmt.Fprintln(os.Stderr, "shelly config doesn't exist. Please run 'shelly init'")
			os.Exit(1)
		}

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

		cmds[args[0]] = cmdstore.Command{
			Name: args[0],
			Description: args[1],
			Action: args[2],
		}

		if err := cmdstore.Save(config.CmdsPath, cmds); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		os.Exit(0)
	},
}
