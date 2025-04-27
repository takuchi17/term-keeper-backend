package main

import (
	"log"

	"github.com/takuchi17/term-keeper/pkg/logger"
)

func main() {
	// set up logger
	logFile, err := logger.SetUpLogger()

	if logFile == nil || err != nil {
		log.Fatal("Failed to open log file: ", err)
	}
	defer logFile.Close()
}
