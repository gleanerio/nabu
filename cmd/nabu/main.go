package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/gleanerio/nabu/pkg/cli"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example

	LOG_FILE := "nabu.log" // log to custom file
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
		return
	}

	log.SetOutput(logFile)

	// Only log the warning severity or above.
	//log.SetLevel(log.WarnLevel)
}

func main() {
	log.Println("calling cli")

	cli.Execute()
}
