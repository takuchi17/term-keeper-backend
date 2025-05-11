package controllers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/takuchi17/term-keeper/api"
	"github.com/takuchi17/term-keeper/app/models"
	"github.com/takuchi17/term-keeper/pkg/jwt"
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

func (h *UserHandeler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		slog.Warn("Invalid request method", "method", r.Method)
		errorResponse := api.ErrorResponse{Message: "Invalid request method"}
		errorJSON, _ := json.Marshal(errorResponse)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, string(errorJSON), http.StatusMethodNotAllowed)
		return
	}

	var ReqestBody api.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&ReqestBody); err != nil {
		slog.Warn("Failed to decode request body", "err", err)
		errorResponse := api.ErrorResponse{Message: "Invalid request body"}
		errorJSON, _ := json.Marshal(errorResponse)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, string(errorJSON), http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByEmail(models.Email(ReqestBody.Email))
	if err != nil {
		slog.Warn("Failed to get user by email", "err", err)
		errorResponse := api.ErrorResponse{Message: "Failed to get user by email"}
		errorJSON, _ := json.Marshal(errorResponse)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, string(errorJSON), http.StatusInternalServerError)
	}

	if err := models.IsSamePassword(user.Password, models.Password(ReqestBody.Password)); err != nil {
		slog.Warn("Failed to check password", "err", err)
		errorResponse := api.ErrorResponse{Message: "Failed to check password"}
		errorJSON, _ := json.Marshal(errorResponse)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, string(errorJSON), http.StatusInternalServerError)
		return
	}

	token, err := jwt.GenerateToken(user.ID, user.Name)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(api.UserLoginResponse{Token: &token})
}
