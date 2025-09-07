package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator - это наш глобальный экземпляр валидатора
var validate = validator.New()

// ValidationError представляет собой более читаемую ошибку валидации
type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// ValidateStruct выполняет валидацию для любой переданной структуры
func ValidateStruct(s interface{}) []*ValidationError {
	var errors []*ValidationError

	err := validate.Struct(s)
	if err != nil {
		// Перебираем ошибки, чтобы создать кастомный ответ
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidationError
			element.Field = strings.ToLower(err.Field()) // Приводим имя поля к нижнему регистру
			element.Error = fmt.Sprintf("failed on the '%s' tag", err.Tag())
			errors = append(errors, &element)
		}
	}

	return errors
}
