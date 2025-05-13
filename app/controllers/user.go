package controllers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/takuchi17/term-keeper/api"
	"github.com/takuchi17/term-keeper/app/models"
	"github.com/takuchi17/term-keeper/pkg/http_checker"
	"github.com/takuchi17/term-keeper/pkg/jwt"
)

type UserHandeler struct {
	DB models.SQLExecutor
}

func (h *UserHandeler) Create(w http.ResponseWriter, r *http.Request) {
	var requestBody api.CreateUserJSONRequestBody
	if err := http_checker.CheckRequest(w, r, &requestBody); err != nil {
		slog.Warn("Failed to check request", "err", err)
		return
	}

	err := models.CreateUser(
		h.DB,
		models.UserName(requestBody.Username),
		models.Email(requestBody.Email),
		models.Password(requestBody.Password),
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
	var requestBody api.UserLoginRequest
	if err := http_checker.CheckRequest(w, r, &requestBody); err != nil {
		slog.Warn("Failed to check request", "err", err)
		return
	}

	user, err := models.GetUserByEmail(h.DB, models.Email(requestBody.Email))
	if err != nil {
		slog.Warn("Failed to get user by email", "err", err)
		errorResponse := api.ErrorResponse{Message: "Failed to get user by email"}
		errorJSON, _ := json.Marshal(errorResponse)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, string(errorJSON), http.StatusInternalServerError)
	}

	if err := models.IsSamePassword(h.DB, user.Password, models.Password(requestBody.Password)); err != nil {
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
