package simulationBatch

import (
	"database/sql"
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/cmd"
	"github.com/jelito/money-maker/app/registry"
	"github.com/jelito/money-maker/app/repository/price"
	"github.com/jelito/money-maker/app/runner/simulationBatch"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Service struct {
}

func (s *Service) Run(configPath *string) {
	config := s.loadConfig(configPath)
	reg := s.createRegistry(config)

	reg.GetByName("app/runner/simulationBatch").(*simulationBatch.Service).Run()
}

func (s *Service) loadConfig(path *string) *config {
	var c config

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

func (s *Service) createRegistry(c *config) *registry.Registry {
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

	cmd.AddDefaultClasses(reg)

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

func (s *Service) createWriter(c *config) *app.Writer {
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
