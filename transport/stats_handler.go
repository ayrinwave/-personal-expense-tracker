package transport

import (
	"Personal-expense-tracking-system/service"
	"encoding/json"
	"net/http"
	"time"
)

// StatsHandler обрабатывает HTTP-запросы, связанные со статистикой
type StatsHandler struct {
	service *service.StatsService
}

// NewStatsHandler создает новый экземпляр StatsHandler
func NewStatsHandler(s *service.StatsService) *StatsHandler {
	return &StatsHandler{service: s}
}

// GetSummary обрабатывает запрос на получение сводки по расходам
func (h *StatsHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		http.Error(w, "could not get user ID from context", http.StatusInternalServerError)
		return
	}

	// Получаем параметры даты из URL
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	// Устанавливаем значения по умолчанию, если параметры не заданы (последние 30 дней)
	to := time.Now()
	from := to.AddDate(0, -1, 0) // Месяц назад

	var err error
	if fromStr != "" {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			http.Error(w, "invalid 'from' date format, use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
	}
	if toStr != "" {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			http.Error(w, "invalid 'to' date format, use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
	}

	// Вызываем сервис для получения сводки
	summary, err := h.service.GetExpenseSummary(r.Context(), userID, from, to)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
