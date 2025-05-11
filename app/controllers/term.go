package controllers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/takuchi17/term-keeper/api"
	"github.com/takuchi17/term-keeper/app/models"
	"github.com/takuchi17/term-keeper/middleware"

	"github.com/takuchi17/term-keeper/pkg/http_checker"
	"github.com/takuchi17/term-keeper/pkg/util"
)

type TermHandler struct {
	DB models.SQLExecutor
}

func (h *TermHandler) Create(w http.ResponseWriter, r *http.Request) {
	var requestBody api.CreateTermJSONRequestBody
	if err := http_checker.CheckRequest(w, r, &requestBody); err != nil {
		slog.Warn("Failed to check request", "err", err)
		return
	}

	userId, ok := middleware.GetUserID(r.Context())
	if !ok {
		slog.Warn("Failed to get user ID from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var categoryIds []models.CategoryId
	for _, categoryId := range *requestBody.CategoryIds {
		categoryIds = append(categoryIds, models.CategoryId(categoryId))
	}

	createdTerm, err := models.CreateTerm(
		h.DB,
		models.TermUserId(userId),
		models.TermName(requestBody.Name),
		models.TermDescription(*requestBody.Description),
		categoryIds,
	)

	if err != nil {
		slog.Error("Failed to create term", "err", err)
		errorResponse := api.ErrorResponse{Message: "Failed to create term"}
		errorJSON, _ := json.Marshal(errorResponse)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, string(errorJSON), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(api.TermResponse{
		Id:          nil,
		Name:        util.Ptr(string(createdTerm.Name)),
		Description: nil,
		CreatedAt:   nil,
		UpdatedAt:   nil,
		Categories:  nil, // Categories are not included in the response
	})
	slog.Info("Term created successfully", "termId", createdTerm.ID)
	w.WriteHeader(http.StatusCreated)
}
