package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/photoshovel/dto"
	"github.com/photoshovel/util"
	"gopkg.in/sconf/ini.v0"
	"gopkg.in/sconf/sconf.v0"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const VERSION = "0.1"
const LOG_FLAG = log.Ldate | log.Ltime | log.Lshortfile
const LOG_DIRECTORY_MODE = os.FileMode(0700)
const LOG_FILE_MODE = os.FileMode(0640)

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

func CreateLogWriter(path string) (logWriter *bufio.Writer) {
	writer, err := util.CreateOrReopenFileForAppending(path, LOG_FILE_MODE)
	
	if err != nil {
		panic(err)
	}
	
	logWriter = bufio.NewWriter(writer)
	
	return
}

func LogInit(config *dto.PhotoShovelConfig) (loggers map[string]*log.Logger, logWriters []*bufio.Writer) {
	var traceHandle, infoHandle, warningHandle, errorHandle *bufio.Writer
	loggers = make(map[string]*log.Logger)
	logWriters = make([]*bufio.Writer, 0, 4)
	
	if config == nil {
		traceHandle = bufio.NewWriter(ioutil.Discard)
		infoHandle = bufio.NewWriter(os.Stdout)
		warningHandle = bufio.NewWriter(os.Stdout)
		errorHandle = bufio.NewWriter(os.Stderr)
	} else {
		// Create the logs directory if it does not already exist.
		err := os.MkdirAll(config.LogDirectory, LOG_DIRECTORY_MODE)
		if err != nil {
			panic(err)
		}
		
		infoHandle = CreateLogWriter(config.LogDirectory + "/photoshovel.info.log")
		warningHandle = CreateLogWriter(config.LogDirectory + "/photoshovel.warning.log")
		errorHandle = CreateLogWriter(config.LogDirectory + "/photoshovel.error.log")
		logWriters = append(logWriters, infoHandle, warningHandle, errorHandle)
		
		if config.TraceLevelLogging {
			traceHandle = CreateLogWriter(config.LogDirectory + "/photoshovel.trace.log")
			logWriters = append(logWriters, traceHandle)
		}
	}
	
	Trace = log.New(traceHandle, "TRACE:", LOG_FLAG)
	loggers["Trace"] = Trace
	Info = log.New(infoHandle, "INFO:", LOG_FLAG)
	loggers["Info"] = Info
	Warning = log.New(warningHandle, "WARNING:", LOG_FLAG)
	loggers["Warning"] = Warning
	Error = log.New(errorHandle, "ERROR:", LOG_FLAG)
	loggers["Error"] = Error
	
	if config == nil {
		Info.Println("No configuration file specified YET.  Continuing to log to console.")
	} else {
		fmt.Printf("Configuration file specified.  Logging to directory:%q.\n", 
			config.LogDirectory)
	}
	
	return
}

func main() {
	// Haven't read config file yet, so use default logging.
	LogInit(nil)
	Info.Println("Welcome to PhotoShovel.")

	configOptPtr := flag.String(
		"config",
		"./photoshovel.config",
		"Path to PhotoShovel configuration file.  Default is './photoshovel.config'.")
	
	flag.Parse()
	Info.Println("Loading config file from: ", *configOptPtr)

	config, err := LoadConfig(configOptPtr)
	if err != nil {
		panicString := fmt.Sprintf(
			"Unable to read config file %q due to %q.  Aborting PhotoShovel.",
			*configOptPtr, err)
		panic(panicString)
	}

	loggers, logWriters := LogInit(config)
	
	_ = loggers
	
	logFlushTicker := time.NewTicker(time.Second)
	go func() {
		for _ = range logFlushTicker.C {
			Trace.Println("Flushing logs.")
			for _, logWriter := range logWriters {
				// Flush logs at least once per second to avoid foul odors.
				logWriter.Flush()
			}
		}
	}()

	Info.Printf("Started PhotoShovel, v.%s.\n", VERSION)
	Info.Printf("Loaded configuration: %q\n", config)

	for _, logWriter := range logWriters {
		// Have logs written to disk before exiting.
		defer logWriter.Flush()
	}
}
