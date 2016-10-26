package main

import (
	"flag"
	"fmt"
	"log"
	//_ "net/http/pprof"
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

	config := flag.String("c", "/etc/dbshield.yml", "config file")
	//cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	listPatterns := flag.Bool("l", false, "get list of captured patterns")
	checkConfig := flag.Bool("k", false, "show parsed config and exit")
	showVersion := flag.Bool("version", false, "show version")
	showHelp := flag.Bool("h", false, "show help")

	//Parsing command line arguments
	flag.Parse()

	// if *cpuprofile != "" {
	// 	f, err := os.Create(*cpuprofile)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	pprof.StartCPUProfile(f)
	// }

	if *showHelp {
		usage(true)
		return
	}

	if *showVersion {
		usage(false)
		return
	}

	if err := dbshield.SetConfigFile(*config); err != nil {
		fmt.Println(err)
		return
	}

	if *listPatterns {
		dbshield.Patterns()
		return
	}

	if *checkConfig {
		err := dbshield.Check()
		log.Println(err)
		return
	}

	log.Println(dbshield.Start())
}
