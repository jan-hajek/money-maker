package cmd

import (
	"database/sql"
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/mailer"
	"github.com/jelito/money-maker/app/registry"
	"github.com/jelito/money-maker/app/repository/position"
	"github.com/jelito/money-maker/app/repository/price"
	"github.com/jelito/money-maker/app/repository/strategy"
	"github.com/jelito/money-maker/app/repository/title"
	"github.com/jelito/money-maker/app/repository/trade"
	"github.com/jelito/money-maker/app/runner/run"
	appTrade "github.com/jelito/money-maker/app/trade"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func init() {
	var cfgFile string

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "start trades for watching",
		Run: func(cmd *cobra.Command, args []string) {
			runCmd := &RunCmd{}
			config := runCmd.loadConfig(&cfgFile)
			reg := runCmd.createRegistry(config)

			reg.GetByName("app/runner/run").(*run.Service).Run()
		},
	}

	runCmd.Flags().StringVar(&cfgFile, "config", "./config.yml", "path to config file")

	rootCmd.AddCommand(runCmd)
}

type RunCmd struct {
}

func (s *RunCmd) loadConfig(path *string) *runConfig {
	var c runConfig

	yamlFile, err := ioutil.ReadFile(*path)
	if err != nil {
		log.Fatalf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return &c
}

func (s *RunCmd) createRegistry(c *runConfig) *registry.Registry {
	reg := registry.Create()

	l := logrus.New()

	l.Hooks.Add(lfshook.NewHook(
		"./data/syslog.log",
		&logrus.JSONFormatter{},
	))

	db, err := sql.Open("mysql", c.Db)
	if err != nil {
		log.Fatal(err)
	}
	reg.Add("db", db)

	AddDefaultClasses(reg)

	reg.Add("app/trade", &appTrade.Factory{
		PositionRepository: reg.GetByName("app/repository/position").(*position.Service),
		Log:                l,
	})

	reg.Add("app/mailer", mailer.Create(
		c.Mail.Enabled,
		c.Mail.Addr,
		c.Mail.From,
		c.Mail.Pass,
		c.Mail.To,
	))

	reg.Add("app/writer", s.createWriter(c))

	reg.Add("app/runner/run", &run.Service{
		Registry:              reg,
		StrategyRepository:    reg.GetByName("app/repository/strategy").(*strategy.Service),
		TradeRepository:       reg.GetByName("app/repository/trade").(*trade.Service),
		PriceRepository:       reg.GetByName("app/repository/price").(*price.Service),
		TitleRepository:       reg.GetByName("app/repository/title").(*title.Service),
		TradeFactory:          reg.GetByName("app/trade").(*appTrade.Factory),
		Log:                   l,
		Writer:                reg.GetByName("app/writer").(*app.Writer),
		Mailer:                reg.GetByName("app/mailer").(*mailer.Service),
		DownloadMissingPrices: c.DownloadMissingPrices,
	})

	return reg
}

func (s *RunCmd) createWriter(c *runConfig) *app.Writer {
	var outputs []app.WriterOutput

	if c.Writer.Outputs.Stdout.Enabled {
		outputs = append(outputs, &app.StdOutWriterOutput{
			DateFormat: c.Writer.ParseFormat,
		})
	}

	if c.Writer.Outputs.Csv.Enabled {
		outputs = append(outputs, &app.CsvWriterOutput{
			File:       c.Writer.Outputs.Csv.Path,
			DateFormat: c.Writer.ParseFormat,
		})
	}

	return &app.Writer{
		Outputs: outputs,
	}
}
