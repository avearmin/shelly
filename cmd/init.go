package cmd

import (
	"fmt"
	"os"

	"github.com/avearmin/shelly/internal/cmdstore"
	"github.com/avearmin/shelly/internal/configstore"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initalize shelly",
	Long: `Initalize shelly by creating .shelly in your home dir, and
		placing commands.json in it.`,
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "HOME env variable not set")
			os.Exit(1)
		}

		if _, err := os.Stat(homeDir + "/.config/shelly/"); err != nil {
			os.Mkdir(homeDir+"/.config/shelly/", 0755)
		}

		if _, err := os.Stat(cmdstore.GetDefaultPath()); err != nil {
			os.Create(cmdstore.GetDefaultPath())
			cmds := make(map[string]cmdstore.Command)
			if err := cmdstore.Save(cmdstore.GetDefaultPath(), cmds); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		if !configstore.Exists() {
			configstore.Create()
			config := configstore.Config{CmdsPath: cmdstore.GetDefaultPath()}
			if err := configstore.Save(config); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		fmt.Println("shelly has been successfully initalized!")
		os.Exit(0)
	},
}
