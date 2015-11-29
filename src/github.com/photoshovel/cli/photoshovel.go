package main

import (
	"flag"
	"fmt"
	"github.com/photoshovel/dto"
	"github.com/photoshovel/util/logging"
	"gopkg.in/sconf/ini.v0"
	"gopkg.in/sconf/sconf.v0"
	"log"
)

const VERSION = "0.1"

var (
	Trace	*log.Logger
	Info	*log.Logger
	Warning	*log.Logger
	Error   *log.Logger
)

func LoadConfig(configFilePath *string) (*dto.PhotoShovelConfig, error) {
	if *configFilePath == "" {
		return nil, nil
	}
	
	config := struct{ PhotoShovel dto.PhotoShovelConfig }{}
	
	err := sconf.Must(&config).Read(ini.File(*configFilePath))
	
	return &config.PhotoShovel, err
}

func main() {
	// Haven't read config file yet, so use default logging.
	loggers, logWriters := logging.LogInit(nil)
	loggers["Info"].Println("Welcome to PhotoShovel.")

	configOptPtr := flag.String(
		"config",
		"./photoshovel.config",
		"Path to PhotoShovel configuration file.  Default is './photoshovel.config'.")
	
	flag.Parse()
	loggers["Info"].Println("Loading config file from: ", *configOptPtr)

	config, err := LoadConfig(configOptPtr)
	if err != nil {
		panicString := fmt.Sprintf(
			"Unable to read config file %q due to %q.  Aborting PhotoShovel.",
			*configOptPtr, err)
		panic(panicString)
	}

	loggers, logWriters = logging.LogInit(config)
	logging.CreateLogFlusher(logWriters, loggers)
	loggers["Info"].Printf("Started PhotoShovel, v.%s.\n", VERSION)
	loggers["Info"].Printf("Loaded configuration: %q\n", config)

	for _, logWriter := range logWriters {
		// Have logs written to disk before exiting.
		defer logWriter.Flush()
	}
}
