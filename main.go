package main

import (
	"log"

	"github.com/takuchi17/term-keeper/app/models"
	"github.com/takuchi17/term-keeper/pkg/logger"
)

func main() {
	// setup logger
	logFile, err := logger.SetUpLogger()
	if logFile == nil || err != nil {
		log.Fatal("Failed to open log file: ", err)
	}
	defer logFile.Close()

	// setup database
	err = models.CreateNewDBConnector(models.InstanceMySQL)
	if err != nil {
		log.Fatal("Failed to connect database: ", err)
	}
	defer models.DB.Close()
}
