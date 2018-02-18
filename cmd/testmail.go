package cmd

import (
	"github.com/jelito/money-maker/src/mailer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	var cfgFile string

	runCmd := &cobra.Command{
		Use:   "test-mail",
		Short: "send test email",
		Run: func(cmd *cobra.Command, args []string) {
			reg := createRegistry(loadConfig(&cfgFile))
			m := reg.GetByName("app/mailer").(*mailer.Service)
			l := reg.GetByName("log").(*logrus.Logger)

			m.ForceEnable()

			l.Info("sending email")
			err := m.Send("test email", "working")
			if err != nil {
				l.Error(err)
			} else {
				l.Info("email sent")
			}
		},
	}

	runCmd.Flags().StringVar(&cfgFile, "config", "./config.yml", "path to config file")

	rootCmd.AddCommand(runCmd)
}
