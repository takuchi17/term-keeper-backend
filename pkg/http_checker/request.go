package http_checker

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/takuchi17/term-keeper/api"
)

func CheckRequest(w http.ResponseWriter, r *http.Request, requestBodyBuf any) error {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		slog.Warn("Invalid request method", "method", r.Method)
		errorResponse := api.ErrorResponse{Message: "Invalid request method"}
		errorJSON, _ := json.Marshal(errorResponse)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, string(errorJSON), http.StatusMethodNotAllowed)
		return errors.New("invalid request method")
	}

	// Check if the content type is application/json
	if err := json.NewDecoder(r.Body).Decode(requestBodyBuf); err != nil {
		slog.Warn("Failed to decode request body", "err", err)
		errorResponse := api.ErrorResponse{Message: "Invalid request body"}
		errorJSON, _ := json.Marshal(errorResponse)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, string(errorJSON), http.StatusBadRequest)
		return errors.New("invalid request body")
	}

	return nil
}
