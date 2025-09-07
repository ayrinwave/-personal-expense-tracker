package rep

import (
	"Personal-expense-tracking-system/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ExpenseRepo отвечает за взаимодействие с таблицей expenses в БД
type ExpenseRepo struct {
	db *pgxpool.Pool
}

// NewExpenseRepo создает новый экземпляр ExpenseRepo
func NewExpenseRepo(db *pgxpool.Pool) *ExpenseRepo {
	return &ExpenseRepo{db: db}
}

// Create сохраняет новый расход в базе данных
func (r *ExpenseRepo) Create(ctx context.Context, expense *models.Expense) error {
	query := `
		INSERT INTO expenses (user_id, category_id, amount, note)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	return r.db.QueryRow(ctx, query, expense.UserID, expense.CategoryID, expense.Amount, expense.Note).Scan(&expense.ID, &expense.CreatedAt)
}

// GetAllByUserID возвращает все расходы для указанного пользователя
func (r *ExpenseRepo) GetAllByUserID(ctx context.Context, userID int) ([]*models.Expense, error) {
	query := `
		SELECT id, user_id, category_id, amount, note, created_at
		FROM expenses
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []*models.Expense
	for rows.Next() {
		var e models.Expense
		if err := rows.Scan(&e.ID, &e.UserID, &e.CategoryID, &e.Amount, &e.Note, &e.CreatedAt); err != nil {
			return nil, err
		}
		expenses = append(expenses, &e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}

// GetByID возвращает один расход по его ID
func (r *ExpenseRepo) GetByID(ctx context.Context, expenseID int) (*models.Expense, error) {
	var e models.Expense
	query := `
		SELECT id, user_id, category_id, amount, note, created_at
		FROM expenses
		WHERE id = $1
	`
	err := r.db.QueryRow(ctx, query, expenseID).Scan(&e.ID, &e.UserID, &e.CategoryID, &e.Amount, &e.Note, &e.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// Update обновляет существующий расход в базе данных
func (r *ExpenseRepo) Update(ctx context.Context, expense *models.Expense) error {
	query := `
		UPDATE expenses
		SET category_id = $1, amount = $2, note = $3
		WHERE id = $4
	`
	_, err := r.db.Exec(ctx, query, expense.CategoryID, expense.Amount, expense.Note, expense.ID)
	return err
}

// Delete удаляет расход из базы данных по его ID
func (r *ExpenseRepo) Delete(ctx context.Context, expenseID int) error {
	query := `DELETE FROM expenses WHERE id = $1`
	_, err := r.db.Exec(ctx, query, expenseID)
	return err
}
