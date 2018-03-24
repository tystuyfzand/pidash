package main

import (
	"github.com/tystuyfzand/pidash/dashboard"
	"github.com/tystuyfzand/pidash/dashboard/config"
	"github.com/tystuyfzand/pidash/dashboard/assets"
	"flag"
	"log"
	"fmt"
)

const Version = "1.0.4"

var (
	flagConfig = flag.String("config", "dashboard.conf", "Config file path")
	flagData = flag.String("data", "./", "Data directory")
	flagDebug = flag.Bool("debug", false, "Enable debugging")
	flagVersion = flag.Bool("version", false, "Show version")
)

func main() {
	flag.Parse()

	if *flagVersion {
		fmt.Println("PiDash version", Version)
		return
	}

	if err := config.Load(*flagConfig); err != nil {
		log.Fatalln("Unable to load configuration:", err)
	}

	config.Debug = *flagDebug

	config.DataDirectory = *flagData

	config.HasLessCompiler = assets.HasLessCompiler()

	dashboard.Dashboard()

	dashboard.Serve()
}