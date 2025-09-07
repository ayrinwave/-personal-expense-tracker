package transport

import (
	"Personal-expense-tracking-system/service"
	"Personal-expense-tracking-system/utils"
	"encoding/json"
	"log"
	"net/http"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// userAuthRequest используется и для регистрации, и для логина
type userAuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// loginResponse — это ответ при успешном входе
type loginResponse struct {
	Token string `json:"token"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req userAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Валидация входных данных
	if errs := utils.ValidateStruct(&req); len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"errors": errs})
		return
	}

	user, err := h.service.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Login обрабатывает запрос на вход пользователя
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req userAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Валидация входных данных
	if errs := utils.ValidateStruct(&req); len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"errors": errs})
		return
	}

	token, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResponse{Token: token})
}
