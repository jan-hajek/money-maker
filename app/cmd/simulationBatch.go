package cmd

import (
	"database/sql"
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/registry"
	"github.com/jelito/money-maker/app/repository/price"
	"github.com/jelito/money-maker/app/runner/simulationBatch"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func init() {
	var cfgFile string

	batchCmd := &cobra.Command{
		Use:   "batch",
		Short: "run batch batchCmd",
		Run: func(cmd *cobra.Command, args []string) {
			batchCmd := SimulationBatchCmd{}
			config := batchCmd.loadConfig(&cfgFile)
			reg := batchCmd.createRegistry(config)
			reg.GetByName("app/runner/simulationBatch").(*simulationBatch.Service).Run()

		},
	}

	batchCmd.Flags().StringVar(&cfgFile, "config", "./config.yml", "path to config file")

	simulationCmd.AddCommand(batchCmd)
}

type SimulationBatchCmd struct {
}

func (s SimulationBatchCmd) loadConfig(path *string) *simulationConfig {
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

func (s SimulationBatchCmd) createRegistry(c *simulationConfig) *registry.Registry {
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

	reg.Add("app/runner/simulationBatch", &simulationBatch.Service{
		PriceRepository: reg.GetByName("app/repository/price").(*price.Service),
		Log:             l,
		Writer:          reg.GetByName("app/writer").(*app.Writer),
		TitleId:         c.TitleId,
		Strategies:      c.Strategies,
		Registry:        reg,
	})

	return reg
}

func (s SimulationBatchCmd) createWriter(c *simulationConfig) *app.Writer {
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
