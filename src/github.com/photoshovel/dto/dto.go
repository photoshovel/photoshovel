package dto

import (
	"fmt"
)

type PhotoShovelConfig struct {
	LogDirectory	string
	TraceLevelLogging bool
	TempFileDirectory string
	NumberOfDownloaders int
	NumberOfUploaders int
}

func (psc *PhotoShovelConfig) String() string {
	return fmt.Sprintf("PhotoShovelConfig:: Log Directory: %q, " +
		"Trace Level Logging: %v, Temp File Directory: %q, Num Downloaders: %v," +
		"Num Uploaders: %v", 
		psc.LogDirectory, psc.TraceLevelLogging, psc.TempFileDirectory, 
		psc.NumberOfDownloaders, psc.NumberOfUploaders)
}

type AmazonCloudDriveConfig struct {
	TargetFolder string
}

type PicasaWebConfig struct {
	TargetAlbumPrefix string
}
