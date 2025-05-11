package main

import (
	"encoding/json"
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/swaggo/swag"

	"github.com/takuchi17/term-keeper/api"
	"github.com/takuchi17/term-keeper/app/controllers"
	"github.com/takuchi17/term-keeper/app/models"
	"github.com/takuchi17/term-keeper/configs"
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
	db, err := models.NewDB(models.InstanceMySQL)
	if err != nil {
		log.Fatal("Failed to connect database: ", err)
	}
	defer db.Close()

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

	userHandler := &controllers.UserHandeler{DB: db}
	http.HandleFunc("/api/v1/signup", userHandler.Create)
	http.HandleFunc("/api/v1/login", userHandler.Login)

	termHandler := &controllers.TermHandler{DB: db}
	http.HandleFunc("/api/v1/terms", termHandler.Create)

	log.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
