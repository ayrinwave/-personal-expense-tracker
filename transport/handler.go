package transport

import (
	"Personal-expense-tracking-system/service"
	"Personal-expense-tracking-system/utils"
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
	if err := utils.DecodeJSON(r.Body, &req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Валидация входных данных
	if errs := utils.ValidateStruct(&req); len(errs) > 0 {
		RespondWithJSON(w, http.StatusBadRequest, map[string]interface{}{"errors": errs})
		return
	}

	user, err := h.service.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		log.Println(err)
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusCreated, user)
}

// Login обрабатывает запрос на вход пользователя
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req userAuthRequest
	if err := utils.DecodeJSON(r.Body, &req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Валидация входных данных
	if errs := utils.ValidateStruct(&req); len(errs) > 0 {
		RespondWithJSON(w, http.StatusBadRequest, map[string]interface{}{"errors": errs})
		return
	}

	token, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		log.Println(err)
		RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	RespondWithJSON(w, http.StatusOK, loginResponse{Token: token})
}
