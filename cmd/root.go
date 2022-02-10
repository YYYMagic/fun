package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "funny-snake",
	Short: "funny-snake is a very fun game that you can play when extremely boared",
	Run: func(cmd *cobra.Command, args []string) {
	// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

