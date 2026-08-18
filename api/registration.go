package api

import (
	"errors"
	"net/http"
	"notesapp/apperrors"
	"notesapp/service"
)

type RegistrationHandler struct {
	registrationService *service.RegistrationService
}

func NewRegistrationHandler(registrationService *service.RegistrationService) *RegistrationHandler {
	return &RegistrationHandler{
		registrationService: registrationService,
	}
}

func (rh RegistrationHandler) Registrate(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")

	err := rh.registrationService.Registration(login, password)
	if errors.Is(err, apperrors.ErrUserAlreadyExists) {
		w.WriteHeader(http.StatusConflict)
		return
	}
	if errors.Is(err, apperrors.ErrNoSuchAuthData) {
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
