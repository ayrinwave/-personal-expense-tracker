package service

import (
	"Personal-expense-tracking-system/models"
	"Personal-expense-tracking-system/rep"
	"context"
	"errors"
)

// ErrForbidden используется, когда пользователь пытается выполнить действие, на которое у него нет прав
var ErrForbidden = errors.New("forbidden")

// ExpenseService содержит бизнес-логику для работы с расходами
type ExpenseService struct {
	repo *rep.ExpenseRepo
}

// NewExpenseService создает новый экземпляр ExpenseService
func NewExpenseService(repo *rep.ExpenseRepo) *ExpenseService {
	return &ExpenseService{repo: repo}
}

// CreateExpense создает новый расход для указанного пользователя
func (s *ExpenseService) CreateExpense(ctx context.Context, userID int, categoryID int, amount float64, note string) (*models.Expense, error) {
	expense := &models.Expense{
		UserID:     userID,
		CategoryID: categoryID,
		Amount:     amount,
		Note:       note,
	}
	if err := s.repo.Create(ctx, expense); err != nil {
		return nil, err
	}
	return expense, nil
}

// GetExpensesByUserID возвращает страницу расходов для указанного пользователя
func (s *ExpenseService) GetExpensesByUserID(ctx context.Context, userID int, page int, limit int) ([]*models.Expense, error) {
	// Вычисляем смещение для SQL-запроса
	offset := (page - 1) * limit
	return s.repo.GetAllByUserID(ctx, userID, limit, offset)
}

// UpdateExpenseRequest содержит поля, которые можно обновить. Указатели используются для частичного обновления.
type UpdateExpenseRequest struct {
	CategoryID *int     `json:"category_id"`
	Amount     *float64 `json:"amount"`
	Note       *string  `json:"note"`
}

// UpdateExpense обновляет существующий расход
func (s *ExpenseService) UpdateExpense(ctx context.Context, userID, expenseID int, req UpdateExpenseRequest) (*models.Expense, error) {
	existingExpense, err := s.repo.GetByID(ctx, expenseID)
	if err != nil {
		return nil, err
	}

	if existingExpense.UserID != userID {
		return nil, ErrForbidden
	}

	if req.CategoryID != nil {
		existingExpense.CategoryID = *req.CategoryID
	}
	if req.Amount != nil {
		existingExpense.Amount = *req.Amount
	}
	if req.Note != nil {
		existingExpense.Note = *req.Note
	}

	if err := s.repo.Update(ctx, existingExpense); err != nil {
		return nil, err
	}

	return existingExpense, nil
}

// DeleteExpense удаляет расход
func (s *ExpenseService) DeleteExpense(ctx context.Context, userID, expenseID int) error {
	// 1. Получаем расход, чтобы проверить владельца
	existingExpense, err := s.repo.GetByID(ctx, expenseID)
	if err != nil {
		return err // Ошибка, если расход не найден
	}

	// 2. Проверяем, что пользователь является владельцем расхода
	if existingExpense.UserID != userID {
		return ErrForbidden // Возвращаем ошибку доступа
	}

	// 3. Если все в порядке, удаляем расход
	return s.repo.Delete(ctx, expenseID)
}
