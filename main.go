package main

import (
	"flag"
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/app/runner"
	"github.com/jelito/money-maker/strategy/jones"
	"github.com/jelito/money-maker/strategy/samson"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func main() {

	runType := flag.String("type", "", "[app, batch]")
	config := flag.String("config", "", "config file")
	flag.Parse()

	var t runner.Config

	yamlFile, err := ioutil.ReadFile(*config)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &t)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	rReg := app.StrategyFactoryRegistry{
		Items: make(map[string]app.StrategyFactory),
	}

	samsonFactory := samson.Factory{}
	rReg.Add(&samsonFactory)

	jonesFactory := jones.Factory{}
	rReg.Add(&jonesFactory)

	switch *runType {
	case "run":
		runner.App{t, &rReg}.Run()
	case "batch":
		runner.App{t, &rReg}.Batch()
	default:
		log.Fatal("unknown param " + *runType)
	}
}
