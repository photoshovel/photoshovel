package logging

import (
	"bufio"
	"fmt"
	"github.com/photoshovel/dto"	
	"github.com/photoshovel/util"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const LOG_DIRECTORY_MODE = os.FileMode(0700)
const LOG_FILE_MODE = os.FileMode(0640)
const LOG_FLAG = log.Ldate | log.Ltime | log.Lshortfile

func CreateLogWriter(path string) (logWriter *bufio.Writer) {
	writer, err := util.CreateOrReopenFileForAppending(path, LOG_FILE_MODE)
	
	if err != nil {
		panic(err)
	}
	
	logWriter = bufio.NewWriter(writer)
	
	return
}

func CreateLogFlusher(logWriters []*bufio.Writer, loggers map[string]*log.Logger) () {
	logFlushTicker := time.NewTicker(time.Second)
	go func() {
		for _ = range logFlushTicker.C {
			loggers["Trace"].Println("Flushing logs.")
			for _, logWriter := range logWriters {
				// Flush logs at least once per second to avoid foul odors.
				logWriter.Flush()
			}
		}
	}()

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
	
	loggers["Trace"] = log.New(traceHandle, "TRACE:", LOG_FLAG)
	loggers["Info"] = log.New(infoHandle, "INFO:", LOG_FLAG)
	loggers["Warning"] = log.New(warningHandle, "WARNING:", LOG_FLAG)
	loggers["Error"] = log.New(errorHandle, "ERROR:", LOG_FLAG)
	
	if config == nil {
		loggers["Info"].Println("No configuration file specified YET.  Continuing to log to console.")
	} else {
		fmt.Printf("Configuration file specified.  Logging to directory:%q.\n", 
			config.LogDirectory)
	}
	
	return
}
