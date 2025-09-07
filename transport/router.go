package transport

import (
	"Personal-expense-tracking-system/service"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// NewRouter теперь принимает оба сервиса и jwtSecret
func NewRouter(userService *service.UserService, expenseService *service.ExpenseService, jwtSecret string) http.Handler {
	r := chi.NewRouter()

	// --- Публичные маршруты ---
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	uh := NewUserHandler(userService)
	r.Post("/user/register", uh.Register)
	r.Post("/user/login", uh.Login)

	// --- Защищенная группа маршрутов ---
	r.Group(func(r chi.Router) {
		// Применяем наш AuthMiddleware ко всем маршрутам в этой группе
		r.Use(AuthMiddleware(jwtSecret))

		// Тестовый маршрут для проверки авторизации
		r.Get("/api/me", func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value(UserIDKey).(int)
			if !ok {
				http.Error(w, "could not get user ID from context", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{"user_id": userID})
		})

		// --- Маршруты для расходов ---
		eh := NewExpenseHandler(expenseService)
		r.Post("/api/expenses", eh.CreateExpense)
		r.Get("/api/expenses", eh.GetExpenses)
		r.Put("/api/expenses/{id}", eh.UpdateExpense)
		r.Delete("/api/expenses/{id}", eh.DeleteExpense) // Добавлено здесь
	})

	return r
}
