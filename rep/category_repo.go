package rep

import (
	"Personal-expense-tracking-system/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CategoryRepo отвечает за взаимодействие с таблицей categories в БД
type CategoryRepo struct {
	db *pgxpool.Pool
}

// NewCategoryRepo создает новый экземпляр CategoryRepo
func NewCategoryRepo(db *pgxpool.Pool) *CategoryRepo {
	return &CategoryRepo{db: db}
}

// Create сохраняет новую категорию в базе данных
func (r *CategoryRepo) Create(ctx context.Context, category *models.Category) error {
	query := `INSERT INTO categories (name) VALUES ($1) RETURNING id`
	return r.db.QueryRow(ctx, query, category.Name).Scan(&category.ID)
}

// GetAll возвращает все категории из базы данных
func (r *CategoryRepo) GetAll(ctx context.Context) ([]*models.Category, error) {
	query := `SELECT id, name FROM categories ORDER BY name ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}

	return categories, rows.Err()
}

// GetByID возвращает одну категорию по ее ID
func (r *CategoryRepo) GetByID(ctx context.Context, id int) (*models.Category, error) {
	var c models.Category
	query := `SELECT id, name FROM categories WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(&c.ID, &c.Name)
	return &c, err
}

// Update обновляет существующую категорию
func (r *CategoryRepo) Update(ctx context.Context, category *models.Category) error {
	query := `UPDATE categories SET name = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, category.Name, category.ID)
	return err
}

// Delete удаляет категорию из базы данных
func (r *CategoryRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
