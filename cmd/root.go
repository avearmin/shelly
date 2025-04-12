package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/avearmin/shelly/internal/cmdstore"
	"github.com/avearmin/shelly/internal/configstore"
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
		if len(selectedCmd) == 0 {
			return
		}

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

		alias := selectedCmd[0]

		shellyCmd, ok := cmds[alias]
		if !ok {
			fmt.Println(alias + " is not a valid alias with shelly.")
			os.Exit(1)
		}
		shellyCmd.LastUsed = time.Now()
		cmds[alias] = shellyCmd
		if err := cmdstore.Save(config.CmdsPath, cmds); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		cmdParts := strings.Fields(selectedCmd[2])

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
