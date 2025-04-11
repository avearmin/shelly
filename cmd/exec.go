package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/avearmin/shelly/internal/cmdstore"
	"github.com/avearmin/shelly/internal/configstore"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(execCmd)
}

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute an shell command using shelly.",
	Long:  "Execute a saved shell command using the assigned alias.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, "'exec' requires args: [ALIAS]")
			os.Exit(1)
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
		
		alias := args[0]

		shellyCmd, ok := cmds[alias]
		if !ok {
			fmt.Println(alias + " is not a valid alias with shelly.")
			os.Exit(1)
		}
		shellyCmd.LastUsed = time.Now()
		cmds[alias] = shellyCmd
		if err := cmdstore.Save(alias, cmds); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		cmdParts := strings.Fields(shellyCmd.Action)

		action := exec.Command(cmdParts[0], cmdParts[1:]...)

		action.Stdin = os.Stdin
		action.Stdout = os.Stdout
		action.Stderr = os.Stderr

		if err := action.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}
