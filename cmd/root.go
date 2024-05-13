package cmd

import (
	"fmt"
	"github.com/avearmin/shelly/internal/tui"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "shelly",
	Short: "shelly manages your shell commands",
	Long: `shelly allows you to save aliases, delete them, and execute 
	their underlying shell commands all from a central location.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := tui.Start(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
