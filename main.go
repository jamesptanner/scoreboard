package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/jamesptanner/scoreboard/src/scoreboard"
)

func main() {
	configFileName := flag.String("config", "", "path to config file to parse")
	outFileName := flag.String("out", "", "output file name.")
	flag.Parse()
	if *configFileName == "" || *outFileName == "" {
		flag.Usage()
		os.Exit(1)
	}
	file, err := os.Open(*configFileName)
	if err != nil {
		log.Printf("Failed to open config file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Failed to read config file: %v\n", err)
		os.Exit(1)
	}
	config := &scoreboard.Config{}

	err = json.Unmarshal(bytes, config)

	if err != nil {
		log.Printf("Failed to unmarshal config file: %v\n", err)
		os.Exit(1)
	}

	scoreboard.RenderBoard(config, outFileName)

}
