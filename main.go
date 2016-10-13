package main

import (
	"flag"
	"log"
	"runtime"

	"github.com/nim4/DBShield/dbshield"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // For Go < 1.5
	//Parsing command line arguments
	config := flag.String("f", "/etc/dbshield.yml", "Config file")
	flag.Parse()

	log.Println(dbshield.Start(*config))
}
