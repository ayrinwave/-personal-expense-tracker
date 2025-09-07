package rep

import (
	"Personal-expense-tracking-system/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUSer(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING id, created_at
	`
	return r.db.QueryRow(ctx, query, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	query := `
		SELECT id, email, password, created_at
		FROM users
		WHERE email=$1
	`
	// ИСПРАВЛЕНО: Теперь сканируем все 4 поля в правильном порядке
	err := r.db.QueryRow(ctx, query, email).Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
