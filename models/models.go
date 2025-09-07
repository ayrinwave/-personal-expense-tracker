package models

import "time"

// User представляет пользователя системы
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Поле Password никогда не отдается клиенту
	CreatedAt time.Time `json:"created_at"`
}

// Expense представляет одну запись о расходе
type Expense struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	CategoryID int       `json:"category_id"`
	Amount     float64   `json:"amount"`
	Note       string    `json:"note"`
	CreatedAt  time.Time `json:"created_at"`
}

// Category представляет одну категорию расходов
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
