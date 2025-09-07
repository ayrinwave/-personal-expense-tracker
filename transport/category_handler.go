package transport

import (
	"Personal-expense-tracking-system/service"
	"Personal-expense-tracking-system/utils"
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
	Name string `json:"name" validate:"required"`
}

// CreateCategory обрабатывает запрос на создание новой категории
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req categoryRequest
	if err := utils.DecodeJSON(r.Body, &req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if errs := utils.ValidateStruct(&req); len(errs) > 0 {
		RespondWithJSON(w, http.StatusBadRequest, map[string]interface{}{"errors": errs})
		return
	}

	category, err := h.service.CreateCategory(r.Context(), req.Name)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, category)
}

// GetAllCategories обрабатывает запрос на получение списка всех категорий
func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAllCategories(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, categories)
}

// UpdateCategory обрабатывает запрос на обновление категории
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := chi.URLParam(r, "id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var req categoryRequest
	if err := utils.DecodeJSON(r.Body, &req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if errs := utils.ValidateStruct(&req); len(errs) > 0 {
		RespondWithJSON(w, http.StatusBadRequest, map[string]interface{}{"errors": errs})
		return
	}

	category, err := h.service.UpdateCategory(r.Context(), categoryID, req.Name)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, category)
}

// DeleteCategory обрабатывает запрос на удаление категории
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := chi.URLParam(r, "id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	if err := h.service.DeleteCategory(r.Context(), categoryID); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
