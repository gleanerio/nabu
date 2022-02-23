package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/gleanerio/nabu/pkg/cli"
)

func init() {
	// Output to stdout instead of the default stderr. Can be any io.Writer, see below for File example
	LOG_FILE := "nabu.log" // log to custom file
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
		return
	}

	log.SetFormatter(&log.JSONFormatter{}) // Log as JSON instead of the default ASCII formatter.
	log.SetReportCaller(true)              // include file name and line number
	log.SetOutput(logFile)

	//log.SetLevel(log.WarnLevel) // Only log the warning severity or above.
}

func main() {
	cli.Execute()
}
