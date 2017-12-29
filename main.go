package main

import (
	"flag"
	"github.com/jelito/money-maker/app"
	"github.com/jelito/money-maker/strategy/jones"
	"github.com/jelito/money-maker/strategy/samson"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	runType := os.Args[1]

	config := flag.String("config", "./config.yml", "config file")
	flag.Parse()

	var t app.Config

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

	samson := samson.Factory{}
	rReg.Add(&samson)

	jones := jones.Factory{}
	rReg.Add(&jones)

	switch runType {
	case "run":
		app.App{t, &rReg}.Run()
	case "batch":
		app.App{t, &rReg}.Batch()
	default:
		panic("unknown param " + runType)
	}
}
