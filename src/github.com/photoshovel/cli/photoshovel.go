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

func LoadConfig(configFilePath *string) (*dto.PhotoShovelConfig, 
										 *dto.AmazonCloudDriveConfig, 
										 *dto.PicasaWebConfig, 
										 error) {
	if *configFilePath == "" {
		return nil, nil, nil, nil
	}
	
	// TODO: Perhaps just return the config struct object?
	config := struct{ PhotoShovel dto.PhotoShovelConfig 
					  AmazonCloudDrive dto.AmazonCloudDriveConfig
					  PicasaWeb dto.PicasaWebConfig }{}
	
	err := sconf.Must(&config).Read(ini.File(*configFilePath))
	
	return &config.PhotoShovel, &config.AmazonCloudDrive, &config.PicasaWeb, err
}

func main() {
	// Haven't read config file yet, so use default logging.
	loggers, logWriters := logging.LogInit(nil)
	loggers["Info"].Println("Welcome to PhotoShovel.")

	configOptPtr := flag.String(
		"config",
		"./photoshovel.config",
		"Path to PhotoShovel configuration file.  Default is './photoshovel.config'.")
	
	sourceOptPtr := flag.String(
		"source",
		"picasaweb",
		"Source of photo migration. [picasaweb|amazonclouddrive]")
	
	targetOptPtr := flag.String(
		"target",
		"amazonclouddrive",
		"Source of photo migration. [picasaweb|amazonclouddrive]")
	
	flag.Parse()
	loggers["Info"].Println("Loading config file from: ", *configOptPtr)

	photoShovelConfig, amazonCloudDriveConfig, picasaWebConfig, err := LoadConfig(configOptPtr)
	_ = amazonCloudDriveConfig
	_ = picasaWebConfig
	if err != nil {
		panicString := fmt.Sprintf(
			"Unable to read config file %q due to %q.  Aborting PhotoShovel.",
			*configOptPtr, err)
		panic(panicString)
	}

	loggers, logWriters = logging.LogInit(photoShovelConfig)
	logging.CreateLogFlusher(logWriters, loggers)
	loggers["Info"].Printf("Started PhotoShovel, v.%s.\n", VERSION)
	loggers["Info"].Printf("Loaded PhotoShovel configuration: %q\n", photoShovelConfig)
	loggers["Info"].Println("Source: ", *sourceOptPtr)
	loggers["Info"].Println("Target: ", *targetOptPtr)

	for _, logWriter := range logWriters {
		// Have logs written to disk before exiting.
		defer logWriter.Flush()
	}
	
	// TODO: Create command line options to specify source and target for 
	// TODO: photo migration.
	
	// TODO: Prompt for photo service authentication if API access keys 
	// TODO: have not been specified in the configuration file.
	
	// TODO: Have Go routine, TaskMaster, running for walking the source photo service.
	// TODO: As it finds files to migrate, writes a migrate job object 
	// TODO: to the downloader channel.
	
	// TODO: Downloader Go routines consume migrate jobs from the downloader 
	// TODO: channel.  For each job, get photo's metadata (album, caption, 
	// TODO: camera info, location, tags, etc) and download the image data 
	// TODO: to a local temporary location.
	
	// TODO: Upon download completion, Download Go routines revise migrate job 
	// TODO: object to contain photo metadata and path to temporary image data 
	// TODO: file.  Migrate job then gets written to uploader channel.
	
	// TODO: Uploader Go routines consume migrate jobs from the uploader 
	// TODO: channel.  For each job, upload the image data+metadata to the target 
	// TODO: photo service.  Delete the temporary image data file upon success.

	// TODO: Once TaskMaster finishes running, write terminate sentinel jobs 
	// TODO: to the downloader and uploader channels equal to the number of 
	// TODO: Downloader and Uploader Go routines in operation.
	
	// TODO: Downloader and Uploader Go routines will consume terminate sentinel 
	// TODO: jobs and stop running.
	
	// TODO: Once all Downloader and Uploader Go routines have exited, end the 
	// TODO: program.
}
