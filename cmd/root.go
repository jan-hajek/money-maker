package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "money-maker",
}

func Execute() {
	rootCmd.Execute()
}
