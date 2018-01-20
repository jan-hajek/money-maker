package cmd

import (
	"database/sql"
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/registry"
	"github.com/jelito/money-maker/app/repository/price"
	"github.com/jelito/money-maker/app/runner/simulationDetail"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func init() {
	var cfgFile string

	detailCmd := &cobra.Command{
		Use:   "detail",
		Short: "show detail",
		Run: func(cmd *cobra.Command, args []string) {
			detailCmd := &SimulationDetailCmd{}
			config := detailCmd.loadConfig(&cfgFile)
			reg := detailCmd.createRegistry(config)
			reg.GetByName("app/runner/simulationDetail").(*simulationDetail.Service).Run()
		},
	}

	detailCmd.Flags().StringVar(&cfgFile, "config", "./config.yml", "path to config file")

	simulationCmd.AddCommand(detailCmd)
}

type SimulationDetailCmd struct {
}

func (s SimulationDetailCmd) loadConfig(path *string) *simulationConfig {
	var c simulationConfig

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

func (s SimulationDetailCmd) createRegistry(c *simulationConfig) *registry.Registry {
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

	reg.Add("app/writer", s.createWriter(c))

	reg.Add("app/runner/simulationDetail", &simulationDetail.Service{
		PriceRepository: reg.GetByName("app/repository/price").(*price.Service),
		Log:             l,
		Writer:          reg.GetByName("app/writer").(*app.Writer),
		TitleId:         c.TitleId,
		Strategies:      c.Strategies,
		Registry:        reg,
	})

	return reg
}

func (s SimulationDetailCmd) createWriter(c *simulationConfig) *app.Writer {
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
