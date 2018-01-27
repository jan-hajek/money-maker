package cmd

import (
	"github.com/spf13/cobra"
	"github.com/takama/daemon"
	"log"
	"os"
)

func init() {
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "install app as service",
		Run: func(cmd *cobra.Command, args []string) {
			srv, err := daemon.New("money-maker", "")
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			srv.Install("run --config=./config.yml")
		},
	}

	rootCmd.AddCommand(installCmd)
}
