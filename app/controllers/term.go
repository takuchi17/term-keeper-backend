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
		Categories:  nil,
	})
	w.WriteHeader(http.StatusCreated)
}

func (h *TermHandler) Get(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r.Context())
	if !ok {
		slog.Warn("Failed to get user ID from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	requestParams := r.URL.Query()
	query := requestParams.Get("query")
	category := requestParams.Get("category")
	sort := requestParams.Get("sort")
	checkedStr := requestParams.Get("checked")

	var termsAndCategories []*models.TermAndCategories
	termsAndCategories, err := models.GetTermsWithCategoriesByUserId(
		h.DB,
		models.TermUserId(userId),
		&query,
		&category,
		&sort,
		&checkedStr,
	)

	if err != nil {
		slog.Error("Failed to get terms", "err", err)
		errorResponse := api.ErrorResponse{Message: "Failed to get terms"}
		errorJSON, _ := json.Marshal(errorResponse)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, string(errorJSON), http.StatusInternalServerError)
		return
	}

	terms := make([]api.TermResponse, len(termsAndCategories))
	for i, termAndCategory := range termsAndCategories {
		categories := make([]api.CategoryResponse, len(termAndCategory.Categories))
		for j, category := range termAndCategory.Categories {
			categories[j] = api.CategoryResponse{
				Id:           util.Ptr(string(category.ID)),
				Name:         util.Ptr(string(category.Name)),
				HexColorCode: util.Ptr(string(category.HexColorCode)),
				CreatedAt:    category.CreatedAt,
				UpdatedAt:    category.UpdatedAt,
			}
		}
		terms[i] = api.TermResponse{
			Id:          util.Ptr(string(termAndCategory.Term.ID)),
			Name:        util.Ptr(string(termAndCategory.Term.Name)),
			Description: util.Ptr(string(termAndCategory.Term.Description)),
			CreatedAt:   termAndCategory.Term.CreatedAt,
			UpdatedAt:   termAndCategory.Term.UpdatedAt,
			Categories:  &categories,
		}
	}

	json.NewEncoder(w).Encode(api.TermListResponse(terms))
}
