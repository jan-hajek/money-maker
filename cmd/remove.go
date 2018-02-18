package cmd

import (
	"github.com/spf13/cobra"
	"github.com/takama/daemon"
	"log"
	"os"
)

func init() {
	removeCmd := &cobra.Command{
		Use:   "remove",
		Short: "remove service",
		Run: func(cmd *cobra.Command, args []string) {
			srv, err := daemon.New("money-maker", "")
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			res, err := srv.Remove()
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			println(res)
		},
	}

	rootCmd.AddCommand(removeCmd)
}
