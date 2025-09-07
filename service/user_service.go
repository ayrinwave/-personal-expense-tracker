package service

import (
	"Personal-expense-tracking-system/models"
	"Personal-expense-tracking-system/rep"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// UserService теперь хранит секрет для генерации JWT
type UserService struct {
	repo      *rep.UserRepo
	jwtSecret string
}

// NewUserService теперь принимает jwtSecret
func NewUserService(repo *rep.UserRepo, jwtSecret string) *UserService {
	return &UserService{repo: repo, jwtSecret: jwtSecret}
}

func (s *UserService) Register(ctx context.Context, email, password string) (*models.User, error) {
	// Проверяем, существует ли уже пользователь
	_, err := s.repo.GetByEmail(ctx, email)
	if err == nil { // Если ошибки нет, значит пользователь найден
		return nil, errors.New("user with this email already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Email:    email,
		Password: string(hashed),
	}

	if err := s.repo.CreateUSer(ctx, user); err != nil {
		return nil, err
	}
	user.Password = "" // Никогда не возвращаем хеш пароля клиенту
	return user, nil
}

// Login реализует логику входа пользователя
func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	// 1. Находим пользователя по email
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		// Если пользователь не найден или произошла ошибка БД, возвращаем общую ошибку
		return "", errors.New("invalid email or password")
	}

	// 2. Сравниваем хеш из БД с паролем, который ввел пользователь
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		// Если пароли не совпадают, возвращаем ту же самую общую ошибку
		return "", errors.New("invalid email or password")
	}

	// 3. Генерируем JWT токен
	claims := jwt.MapClaims{
		"sub": user.ID,                               // ID пользователя
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Срок действия токена - 24 часа
		"iat": time.Now().Unix(),                     // Время создания
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 4. Подписываем токен нашим секретным ключом
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
