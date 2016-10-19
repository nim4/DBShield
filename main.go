package main

import (
	"flag"
	"log"
	_ "net/http/pprof"
	"runtime"

	"github.com/nim4/DBShield/dbshield"
)

func usage(showUsage bool) {
	print("DBShield " + dbshield.Version + "\n")
	if showUsage {
		flag.Usage()
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // For Go < 1.5
	//Parsing command line arguments
	config := flag.String("c", "/etc/dbshield.yml", "Config file")
	checkConfig := flag.Bool("k", false, "Show parsed config and exit")
	showVersion := flag.Bool("version", false, "Show version")
	showHelp := flag.Bool("h", false, "Show help")
	flag.Parse()

	if *showHelp {
		usage(true)
		return
	}

	if *showVersion {
		usage(false)
		return
	}

	if *checkConfig {
		if err := dbshield.Check(*config); err != nil {
			log.Println(err)
		}
		return
	}

	log.Println(dbshield.Start(*config))
}
