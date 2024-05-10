package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use:   "shelly",
	Short: "shelly manages your shell commands",
	Long: `shelly allows you to save aliases, delete them, and execute 
	their underlying shell commands all from a central location.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("There will be a full TUI here soon!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
