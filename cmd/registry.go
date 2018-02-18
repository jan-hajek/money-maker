package cmd

import (
	"database/sql"
	"github.com/Gurpartap/logrus-stack"
	"github.com/jelito/money-maker/app/strategy/jones"
	"github.com/jelito/money-maker/app/strategy/jones2"
	"github.com/jelito/money-maker/app/strategy/samson"
	"github.com/jelito/money-maker/app/title/admiralMarkets"
	"github.com/jelito/money-maker/app/title/plus500"
	"github.com/jelito/money-maker/src/dateInput"
	"github.com/jelito/money-maker/src/mailer"
	"github.com/jelito/money-maker/src/position"
	"github.com/jelito/money-maker/src/registry"
	positionRepo "github.com/jelito/money-maker/src/repository/position"
	"github.com/jelito/money-maker/src/repository/price"
	strategyRepo "github.com/jelito/money-maker/src/repository/strategy"
	"github.com/jelito/money-maker/src/repository/title"
	tradeRepo "github.com/jelito/money-maker/src/repository/trade"
	"github.com/jelito/money-maker/src/runner/run"
	"github.com/jelito/money-maker/src/runner/simulationBatch"
	"github.com/jelito/money-maker/src/runner/simulationDetail"
	"github.com/jelito/money-maker/src/trade"
	"github.com/jelito/money-maker/src/writer"
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

	reg.Add("app/repository/trade", &tradeRepo.Service{Db: db})
	reg.Add("app/repository/title", &title.Service{Db: db})
	reg.Add("app/repository/price", &price.Service{Db: db})
	reg.Add("app/repository/strategy", &strategyRepo.Service{Db: db})
	reg.Add("app/repository/position", &positionRepo.Service{Db: db})

	reg.Add("strategy/samson", &samson.Factory{})
	reg.Add("strategy/jones", &jones.Factory{})
	reg.Add("strategy/jones2", &jones2.Factory{})
	reg.Add("title/admiralMarkets", &admiralMarkets.Factory{})
	reg.Add("title/plus500", &plus500.Factory{})

	reg.Add("app/position/maker", &position.PositionMaker{})

	reg.Add("app/trade", &trade.Factory{
		PositionMaker: reg.GetByName("app/position/maker").(*position.PositionMaker),
		Log:           l,
	})

	reg.Add("app/runner/run", &run.Service{
		Registry:           reg,
		StrategyRepository: reg.GetByName("app/repository/strategy").(*strategyRepo.Service),
		TradeRepository:    reg.GetByName("app/repository/trade").(*tradeRepo.Service),
		PriceRepository:    reg.GetByName("app/repository/price").(*price.Service),
		TitleRepository:    reg.GetByName("app/repository/title").(*title.Service),
		TradeFactory:       reg.GetByName("app/trade").(*trade.Factory),
		Log:                l,
		Writer:             reg.GetByName("app/writer").(*writer.Writer),
		MailBufferFactory: createMailBufferFactory(
			reg.GetByName("app/mailer").(*mailer.Service),
			l,
		),
		DownloadMissingPrices: c.Run.DownloadMissingPrices,
	})

	reg.Add("app/runner/simulationBatch", &simulationBatch.Service{
		Log:             l,
		Writer:          reg.GetByName("app/writer").(*writer.Writer),
		Strategies:      c.Simulation.Strategies,
		Registry:        reg,
		DateInputLoader: createDateInputLoader(c, l, reg),
	})

	reg.Add("app/runner/simulationDetail", &simulationDetail.Service{
		Log:             l,
		Writer:          reg.GetByName("app/writer").(*writer.Writer),
		Strategies:      c.Simulation.Strategies,
		Registry:        reg,
		DateInputLoader: createDateInputLoader(c, l, reg),
	})

	return reg
}

func createLogger() *logrus.Logger {
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	l.Hooks.Add(lfshook.NewHook(
		"./data/syslog.log",
		&logrus.TextFormatter{},
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
		c.Mail.Host,
	)
}

func createMailBufferFactory(service *mailer.Service, l *logrus.Logger) *mailer.BufferFactory {
	return &mailer.BufferFactory{
		Mailer: service,
		Log:    l,
	}

}

func createWriter(c *config) *writer.Writer {
	var outputs []writer.Output

	if c.Writer.Outputs.Stdout.Enabled {
		outputs = append(outputs, &writer.StdOutWriterOutput{
			DateFormat: c.Writer.ParseFormat,
		})
	}

	if c.Writer.Outputs.Csv.Enabled {
		outputs = append(outputs, &writer.CsvWriterOutput{
			File:       c.Writer.Outputs.Csv.Path,
			DateFormat: c.Writer.ParseFormat,
		})
	}

	return &writer.Writer{
		Outputs: outputs,
	}
}

func createDateInputLoader(c *config, l *logrus.Logger, reg *registry.Registry) dateInput.Loader {
	cs := c.Simulation.Source

	if cs.Db.Enabled && cs.Csv.Enabled {
		l.Fatal("enable only one simulation source in config")
	}

	if cs.Db.Enabled {
		return &dateInput.DatabaseLoader{
			TitleId:         cs.Db.TitleId,
			PriceRepository: reg.GetByName("app/repository/price").(*price.Service),
		}
	}

	if cs.Csv.Enabled {
		return &dateInput.CsvLoader{
			InputFilePath:   cs.Csv.FilePath,
			TimeParseFormat: cs.Csv.TimeParseFormat,
		}
	}

	return nil
}
