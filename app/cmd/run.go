package cmd

import (
	"github.com/jelito/money-maker/app/runner/run"
	"github.com/spf13/cobra"
)

func init() {
	var cfgFile string

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "start trades for watching",
		Run: func(cmd *cobra.Command, args []string) {
			reg := createRegistry(loadConfig(&cfgFile))
			reg.GetByName("app/runner/run").(*run.Service).Run()
		},
	}

	runCmd.Flags().StringVar(&cfgFile, "config", "./config.yml", "path to config file")

	rootCmd.AddCommand(runCmd)
}
