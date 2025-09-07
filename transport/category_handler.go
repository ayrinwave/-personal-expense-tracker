package transport

import (
	"Personal-expense-tracking-system/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// CategoryHandler обрабатывает HTTP-запросы, связанные с категориями
type CategoryHandler struct {
	service *service.CategoryService
}

// NewCategoryHandler создает новый экземпляр CategoryHandler
func NewCategoryHandler(s *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
}

// categoryRequest — это структура для парсинга тела запроса на создание/обновление категории
type categoryRequest struct {
	Name string `json:"name"`
}

// CreateCategory обрабатывает запрос на создание новой категории
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req categoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	category, err := h.service.CreateCategory(r.Context(), req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

// GetAllCategories обрабатывает запрос на получение списка всех категорий
func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAllCategories(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// UpdateCategory обрабатывает запрос на обновление категории
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := chi.URLParam(r, "id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		http.Error(w, "invalid category ID", http.StatusBadRequest)
		return
	}

	var req categoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	category, err := h.service.UpdateCategory(r.Context(), categoryID, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// DeleteCategory обрабатывает запрос на удаление категории
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := chi.URLParam(r, "id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		http.Error(w, "invalid category ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteCategory(r.Context(), categoryID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
