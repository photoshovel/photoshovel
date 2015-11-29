package dto

import (
	"fmt"
)

type PhotoShovelConfig struct {
	LogDirectory	string
	TraceLevelLogging bool
}

func (psc *PhotoShovelConfig) String() string {
	return fmt.Sprintf("PhotoShovelConfig:: Log Directory: %q, " +
		"Trace Level Logging: %v", 
		psc.LogDirectory, psc.TraceLevelLogging)
}