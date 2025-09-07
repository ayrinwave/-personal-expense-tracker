package service

import (
	"Personal-expense-tracking-system/rep"
	"context"
	"time"
)

// StatsService содержит бизнес-логику для работы со статистикой
type StatsService struct {
	expenseRepo *rep.ExpenseRepo
}

// NewStatsService создает новый экземпляр StatsService
func NewStatsService(expenseRepo *rep.ExpenseRepo) *StatsService {
	return &StatsService{expenseRepo: expenseRepo}
}

// ExpenseSummary представляет собой сводку по расходам
type ExpenseSummary struct {
	TotalAmount float64   `json:"total_amount"`
	From        time.Time `json:"from"`
	To          time.Time `json:"to"`
}

// GetExpenseSummary возвращает общую сумму расходов пользователя за период
func (s *StatsService) GetExpenseSummary(ctx context.Context, userID int, from, to time.Time) (*ExpenseSummary, error) {
	total, err := s.expenseRepo.GetTotalAmountByUserIDAndDateRange(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}

	summary := &ExpenseSummary{
		TotalAmount: total,
		From:        from,
		To:          to,
	}

	return summary, nil
}
