package cmd

import (
	"github.com/spf13/cobra"
)

var simulationCmd = &cobra.Command{
	Use:   "simulation",
	Short: "start simulation",
}

func init() {
	rootCmd.AddCommand(simulationCmd)
}
