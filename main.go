package main

import (
	"flag"
	"github.com/jelito/money-maker/app/cmd/run"
	"log"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	runType := flag.String("type", "", "[run, app, batch]")
	config := flag.String("config", "", "config file")
	flag.Parse()

	switch *runType {
	case "run":
		s := &run.Service{}
		s.Run(config)
	default:
		log.Fatal("unknown param " + *runType)
	}
}
