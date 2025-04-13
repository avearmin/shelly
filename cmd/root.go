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
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
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

		cmdSlice := make([]cmdstore.Command, len(cmds))
		i := 0
		for _, v := range cmds {
			cmdSlice[i] = v
			i++
		}

		selectedCmd, err := tui.Start(cmdSlice)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		selectedCmd.LastUsed = time.Now()
		cmds[selectedCmd.Name] = selectedCmd
		if err := cmdstore.Save(config.CmdsPath, cmds); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		cmdParts := strings.Fields(selectedCmd.Action)

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
