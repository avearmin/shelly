package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
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

		if _, err := os.Stat(homeDir + "/.config/shelly/commands.json"); err != nil {
			os.Create(homeDir + "/.config/shelly/commands.json")
		}

		fmt.Println("shelly has been successfully initalized!")
		os.Exit(0)
	},
}
