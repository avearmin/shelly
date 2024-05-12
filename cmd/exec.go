package cmd

import (
	"fmt"
	"github.com/avearmin/shelly/internal/cmdstore"
	"github.com/avearmin/shelly/internal/configstore"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
)

func init() {
	rootCmd.AddCommand(execCmd)
}

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute an shell command using shelly.",
	Long:  "Execute a saved shell command using the assigned alias.",
	Run: func(cmd *cobra.Command, args []string) {
		if !configstore.Exists() {
			fmt.Fprintln(os.Stderr, "shelly config doesn't exist. Please run 'shelly init'")
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

		cmdParts := strings.Fields(cmds[args[0]].Action)

		action := exec.Command(cmdParts[0], cmdParts[1:]...)

		action.Stdin = os.Stdin
		action.Stdout = os.Stdout
		action.Stderr = os.Stderr

		if err := action.Start(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := action.Wait(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}
