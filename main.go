package main

import (
	"flag"
	"github.com/jelito/money-maker/app/cmd/run"
	"github.com/jelito/money-maker/app/cmd/simulationBatch"
	"github.com/jelito/money-maker/app/cmd/simulationDetail"
	"log"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	runType := flag.String("type", "", "[run, simulationBatch]")
	config := flag.String("config", "", "config file")
	flag.Parse()

	switch *runType {
	case "run":
		s := &run.Service{}
		s.Run(config)
	case "simulationBatch":
		s := &simulationBatch.Service{}
		s.Run(config)
	case "simulationDetail":
		s := &simulationDetail.Service{}
		s.Run(config)
	default:
		log.Fatal("unknown param " + *runType)
	}
}
