package transport

import (
	"Personal-expense-tracking-system/service"
	"Personal-expense-tracking-system/utils"
	"encoding/json"
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
		http.Error(w, "could not get user ID from context", http.StatusInternalServerError)
		return
	}

	var req createExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Валидация входных данных
	if errs := utils.ValidateStruct(&req); len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"errors": errs})
		return
	}

	expense, err := h.service.CreateExpense(r.Context(), userID, req.CategoryID, req.Amount, req.Note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(expense)
}

// GetExpenses обрабатывает запрос на получение списка расходов пользователя
func (h *ExpenseHandler) GetExpenses(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		http.Error(w, "could not get user ID from context", http.StatusInternalServerError)
		return
	}

	expenses, err := h.service.GetExpensesByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

// UpdateExpense обрабатывает запрос на обновление расхода
func (h *ExpenseHandler) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		http.Error(w, "could not get user ID from context", http.StatusInternalServerError)
		return
	}

	expenseIDStr := chi.URLParam(r, "id")
	expenseID, err := strconv.Atoi(expenseIDStr)
	if err != nil {
		http.Error(w, "invalid expense ID", http.StatusBadRequest)
		return
	}

	var req service.UpdateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Валидация входных данных
	if errs := utils.ValidateStruct(&req); len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"errors": errs})
		return
	}

	updatedExpense, err := h.service.UpdateExpense(r.Context(), userID, expenseID, req)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			http.Error(w, "you don't have permission to update this expense", http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedExpense)
}

// DeleteExpense обрабатывает запрос на удаление расхода
func (h *ExpenseHandler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		http.Error(w, "could not get user ID from context", http.StatusInternalServerError)
		return
	}

	expenseIDStr := chi.URLParam(r, "id")
	expenseID, err := strconv.Atoi(expenseIDStr)
	if err != nil {
		http.Error(w, "invalid expense ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteExpense(r.Context(), userID, expenseID); err != nil {
		if errors.Is(err, service.ErrForbidden) {
			http.Error(w, "you don't have permission to delete this expense", http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
