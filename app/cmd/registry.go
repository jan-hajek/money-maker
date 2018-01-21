package cmd

import (
	"database/sql"
	"github.com/Gurpartap/logrus-stack"
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/mailer"
	"github.com/jelito/money-maker/app/registry"
	"github.com/jelito/money-maker/app/repository/position"
	"github.com/jelito/money-maker/app/repository/price"
	"github.com/jelito/money-maker/app/repository/strategy"
	"github.com/jelito/money-maker/app/repository/title"
	"github.com/jelito/money-maker/app/repository/trade"
	"github.com/jelito/money-maker/app/runner/run"
	"github.com/jelito/money-maker/app/runner/simulationBatch"
	"github.com/jelito/money-maker/app/runner/simulationDetail"
	appTrade "github.com/jelito/money-maker/app/trade"
	"github.com/jelito/money-maker/strategy/jones"
	"github.com/jelito/money-maker/strategy/jones2"
	"github.com/jelito/money-maker/strategy/samson"
	"github.com/jelito/money-maker/title/admiralMarkets"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func loadConfig(path *string) *config {
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

func createRegistry(c *config) *registry.Registry {
	reg := registry.Create()

	l := createLogger()
	reg.Add("log", l)

	db := createDb(c, l)
	reg.Add("db", db)

	reg.Add("app/mailer", createMailer(c))
	reg.Add("app/writer", createWriter(c))

	reg.Add("app/repository/trade", &trade.Service{Db: db})
	reg.Add("app/repository/title", &title.Service{Db: db})
	reg.Add("app/repository/price", &price.Service{Db: db})
	reg.Add("app/repository/strategy", &strategy.Service{Db: db})
	reg.Add("app/repository/position", &position.Service{Db: db})

	reg.Add("strategy/samson", &samson.Factory{})
	reg.Add("strategy/jones", &jones.Factory{})
	reg.Add("strategy/jones2", &jones2.Factory{})
	reg.Add("title/admiralMarkets", &admiralMarkets.Factory{})

	reg.Add("app/trade", &appTrade.Factory{
		PositionRepository: reg.GetByName("app/repository/position").(*position.Service),
		Log:                l,
	})

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
		DownloadMissingPrices: c.Run.DownloadMissingPrices,
	})

	reg.Add("app/runner/simulationBatch", &simulationBatch.Service{
		PriceRepository: reg.GetByName("app/repository/price").(*price.Service),
		Log:             l,
		Writer:          reg.GetByName("app/writer").(*app.Writer),
		TitleId:         c.Simulation.TitleId,
		Strategies:      c.Simulation.Strategies,
		Registry:        reg,
	})

	reg.Add("app/runner/simulationDetail", &simulationDetail.Service{
		PriceRepository: reg.GetByName("app/repository/price").(*price.Service),
		Log:             l,
		Writer:          reg.GetByName("app/writer").(*app.Writer),
		TitleId:         c.Simulation.TitleId,
		Strategies:      c.Simulation.Strategies,
		Registry:        reg,
	})

	return reg
}
func createLogger() *logrus.Logger {
	l := logrus.New()
	l.Hooks.Add(lfshook.NewHook(
		"./data/syslog.log",
		&logrus.JSONFormatter{},
	))

	callerLevels := []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel}
	stackLevels := []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel}

	l.AddHook(logrus_stack.NewHook(callerLevels, stackLevels))

	return l
}

func createDb(c *config, l *logrus.Logger) *sql.DB {
	db, err := sql.Open("mysql", c.Db)
	if err != nil {
		l.Fatal(err)
	}
	return db
}

func createMailer(c *config) *mailer.Service {
	return mailer.Create(
		c.Mail.Enabled,
		c.Mail.Addr,
		c.Mail.From,
		c.Mail.Pass,
		c.Mail.To,
	)
}

func createWriter(c *config) *app.Writer {
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
