package cmd

import (
	"github.com/jelito/money-maker/src/runner/simulationBatch"
	"github.com/jelito/money-maker/src/runner/simulationDetail"
	"github.com/spf13/cobra"
)

func init() {
	var cfgFile string

	simulationCmd := &cobra.Command{
		Use:   "simulation",
		Short: "start simulation",
	}

	batchCmd := &cobra.Command{
		Use:   "batch",
		Short: "run batch batchCmd",
		Run: func(cmd *cobra.Command, args []string) {
			reg := createRegistry(loadConfig(&cfgFile))
			reg.GetByName("app/runner/simulationBatch").(*simulationBatch.Service).Run()

		},
	}
	batchCmd.Flags().StringVar(&cfgFile, "config", "./config.yml", "path to config file")

	detailCmd := &cobra.Command{
		Use:   "detail",
		Short: "show detail",
		Run: func(cmd *cobra.Command, args []string) {
			reg := createRegistry(loadConfig(&cfgFile))
			reg.GetByName("app/runner/simulationDetail").(*simulationDetail.Service).Run()
		},
	}
	detailCmd.Flags().StringVar(&cfgFile, "config", "./config.yml", "path to config file")

	simulationCmd.AddCommand(batchCmd)
	simulationCmd.AddCommand(detailCmd)
	rootCmd.AddCommand(simulationCmd)
}
