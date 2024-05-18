package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/avearmin/shelly/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "shelly",
	Short: "shelly manages your shell commands",
	Long: `shelly allows you to save aliases, delete them, and execute
		their underlying shell commands all from a central location.`,
	Run: func(cmd *cobra.Command, args []string) {
		selectedCmd, err := tui.Start()
		if err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
		if selectedCmd == "" {
			return
		}

		cmdParts := strings.Fields(selectedCmd)

		action := exec.Command(cmdParts[0], cmdParts[1:]...)

		action.Stdin = os.Stdin
		action.Stdout = os.Stdout
		action.Stderr = os.Stderr

		if err := action.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
