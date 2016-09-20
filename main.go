package main

import (
	"flag"
	"log"

	"./dbshield"
)

func main() {
	//Parsing command line arguments
	config := flag.String("f", "/etc/dbshield.yml", "Config file")
	flag.Parse()

	log.Println(dbshield.Start(*config))
}
