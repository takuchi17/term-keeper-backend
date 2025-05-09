package main

import (
	"encoding/json"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/swaggo/swag"
	"github.com/takuchi17/term-keeper/api"
	"github.com/takuchi17/term-keeper/app/models"
	"github.com/takuchi17/term-keeper/configs"
	"github.com/takuchi17/term-keeper/pkg/logger"
	"log"
	"net/http"
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

	swagger, err := api.GetSwagger()
	if err != nil {
		log.Fatal("Failed to generate swagger: ", err)
	}
	if configs.Config.IsDevelopment() {
		swaggerJson, _ := json.Marshal(swagger)
		var SwaggerInfo = &swag.Spec{
			InfoInstanceName: "swagger",
			SwaggerTemplate:  string(swaggerJson),
		}
		swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
		http.Handle("/swagger/", httpSwagger.WrapHandler)
	}

	log.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
