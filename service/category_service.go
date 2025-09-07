package service

import (
	"Personal-expense-tracking-system/models"
	"Personal-expense-tracking-system/rep"
	"context"
)

// CategoryService содержит бизнес-логику для работы с категориями
type CategoryService struct {
	repo *rep.CategoryRepo
}

// NewCategoryService создает новый экземпляр CategoryService
func NewCategoryService(repo *rep.CategoryRepo) *CategoryService {
	return &CategoryService{repo: repo}
}

// CreateCategory создает новую категорию
func (s *CategoryService) CreateCategory(ctx context.Context, name string) (*models.Category, error) {
	category := &models.Category{
		Name: name,
	}
	if err := s.repo.Create(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

// GetAllCategories возвращает все категории
func (s *CategoryService) GetAllCategories(ctx context.Context) ([]*models.Category, error) {
	return s.repo.GetAll(ctx)
}

// UpdateCategory обновляет существующую категорию
func (s *CategoryService) UpdateCategory(ctx context.Context, id int, name string) (*models.Category, error) {
	category := &models.Category{
		ID:   id,
		Name: name,
	}
	if err := s.repo.Update(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

// DeleteCategory удаляет категорию
func (s *CategoryService) DeleteCategory(ctx context.Context, id int) error {
	// Примечание: в реальном приложении здесь стоило бы проверить,
	// не используются ли эта категория в каких-либо расходах, прежде чем удалять.
	// Пока для простоты мы это опускаем.
	return s.repo.Delete(ctx, id)
}
