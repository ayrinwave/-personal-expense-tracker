package transport

import (
	"Personal-expense-tracking-system/service"
	"Personal-expense-tracking-system/utils"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// ExpenseHandler обрабатывает HTTP-запросы, связанные с расходами
type ExpenseHandler struct {
	service *service.ExpenseService
}

// NewExpenseHandler создает новый экземпляр ExpenseHandler
func NewExpenseHandler(s *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{service: s}
}

// createExpenseRequest — это структура для парсинга тела запроса на создание расхода
type createExpenseRequest struct {
	CategoryID int     `json:"category_id" validate:"required,gt=0"`
	Amount     float64 `json:"amount" validate:"required,gt=0"`
	Note       string  `json:"note"`
}

// CreateExpense обрабатывает запрос на создание нового расхода
func (h *ExpenseHandler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		RespondWithError(w, http.StatusInternalServerError, "Could not get user ID from context")
		return
	}

	var req createExpenseRequest
	if err := utils.DecodeJSON(r.Body, &req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if errs := utils.ValidateStruct(&req); len(errs) > 0 {
		RespondWithJSON(w, http.StatusBadRequest, map[string]interface{}{"errors": errs})
		return
	}

	expense, err := h.service.CreateExpense(r.Context(), userID, req.CategoryID, req.Amount, req.Note)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, expense)
}

// GetExpenses обрабатывает запрос на получение списка расходов пользователя с пагинацией
func (h *ExpenseHandler) GetExpenses(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		RespondWithError(w, http.StatusInternalServerError, "Could not get user ID from context")
		return
	}

	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		RespondWithError(w, http.StatusBadRequest, "Invalid page parameter")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "10"
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		RespondWithError(w, http.StatusBadRequest, "Invalid limit parameter")
		return
	}

	expenses, err := h.service.GetExpensesByUserID(r.Context(), userID, page, limit)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, expenses)
}

// UpdateExpense обрабатывает запрос на обновление расхода
func (h *ExpenseHandler) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		RespondWithError(w, http.StatusInternalServerError, "Could not get user ID from context")
		return
	}

	expenseIDStr := chi.URLParam(r, "id")
	expenseID, err := strconv.Atoi(expenseIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid expense ID")
		return
	}

	var req service.UpdateExpenseRequest
	if err := utils.DecodeJSON(r.Body, &req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if errs := utils.ValidateStruct(&req); len(errs) > 0 {
		RespondWithJSON(w, http.StatusBadRequest, map[string]interface{}{"errors": errs})
		return
	}

	updatedExpense, err := h.service.UpdateExpense(r.Context(), userID, expenseID, req)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			RespondWithError(w, http.StatusForbidden, "You don't have permission to update this expense")
		} else {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondWithJSON(w, http.StatusOK, updatedExpense)
}

// DeleteExpense обрабатывает запрос на удаление расхода
func (h *ExpenseHandler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		RespondWithError(w, http.StatusInternalServerError, "Could not get user ID from context")
		return
	}

	expenseIDStr := chi.URLParam(r, "id")
	expenseID, err := strconv.Atoi(expenseIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid expense ID")
		return
	}

	if err := h.service.DeleteExpense(r.Context(), userID, expenseID); err != nil {
		if errors.Is(err, service.ErrForbidden) {
			RespondWithError(w, http.StatusForbidden, "You don't have permission to delete this expense")
		} else {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
