package transport

import (
	"Personal-expense-tracking-system/service"
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
		RespondWithError(w, http.StatusInternalServerError, "Could not get user ID from context")
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
			RespondWithError(w, http.StatusBadRequest, "Invalid 'from' date format, use YYYY-MM-DD")
			return
		}
	}
	if toStr != "" {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid 'to' date format, use YYYY-MM-DD")
			return
		}
	}

	// Вызываем сервис для получения сводки
	summary, err := h.service.GetExpenseSummary(r.Context(), userID, from, to)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, summary)
}
