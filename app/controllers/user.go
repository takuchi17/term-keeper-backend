package controllers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/takuchi17/term-keeper/api"
	"github.com/takuchi17/term-keeper/app/models"
	// "github.com/takuchi17/term-keeper/app/models"
)

type UserHandeler struct{}

func (h *UserHandeler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var ReqestBody api.CreateUserJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&ReqestBody); err != nil {
		slog.Warn("Failed to decode request body", "err", err)
		errorRespose := api.ErrorResponse{Message: "Invalid request body"}
		errorJSON, _ := json.Marshal(errorRespose)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, string(errorJSON), http.StatusBadRequest)
		return
	}

	err := models.CreateUser(
		models.UserName(ReqestBody.Username),
		models.Email(ReqestBody.Email),
		models.Password(ReqestBody.Password),
	)

	if err != nil {
		slog.Error("Failed to create user", "err", err)
		errorResponse := api.ErrorResponse{Message: "Failed to create user"}
		errorJSON, _ := json.Marshal(errorResponse)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, string(errorJSON), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
