package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"notesapp/apperrors"
	"notesapp/dto"
	"notesapp/service"
	"strconv"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	UserService *service.UserService
}

func (ud UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userEntity, err := ud.UserService.GetById(userId)
	if errors.Is(err, apperrors.ErrNoSuchUser) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userResponse := dto.ToUserResponse(userEntity)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userResponse)
}

func (ud UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userEntities, err := ud.UserService.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userResponses := []dto.UserResponse{}
	for _, userEntity := range userEntities {
		userResponses = append(userResponses, dto.ToUserResponse(userEntity))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userResponses)
}

func (ud UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := cookie.Value
	user, err := ud.UserService.GetMe(token)
	if errors.Is(err, apperrors.ErrNoSuchSession) || errors.Is(err, apperrors.ErrSessionExpired) || errors.Is(err, apperrors.ErrNoSuchUser) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userResponse := dto.ToUserResponse(user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userResponse)
}

func NewUserHandler(us *service.UserService) UserHandler {
	return UserHandler{
		UserService: us,
	}
}
