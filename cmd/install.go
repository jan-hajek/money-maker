package cmd

import (
	"github.com/spf13/cobra"
	"github.com/takama/daemon"
	"log"
	"os"
)

func init() {
	var cfgFile string

	installCmd := &cobra.Command{
		Use:   "install",
		Short: "install app as service",
		Run: func(cmd *cobra.Command, args []string) {
			srv, err := daemon.New("money-maker", "")
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			res, err := srv.Install("run --config=" + *&cfgFile)

			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			println(res)
		},
	}

	installCmd.Flags().StringVar(&cfgFile, "config", "./config.yml", "path to config file")

	rootCmd.AddCommand(installCmd)
}
