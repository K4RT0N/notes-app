package api

import (
	"errors"
	"net/http"
	"notesapp/apperrors"
	"notesapp/service"
	"time"
)

type AuthHandler struct {
	sessionService *service.SessionService
}

func NewAuthHandler(sessionService *service.SessionService) *AuthHandler {
	return &AuthHandler{
		sessionService: sessionService,
	}
}

func (ah *AuthHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")

	token, err := ah.sessionService.InitializeSession(login, password)
	if errors.Is(err, apperrors.ErrWrongCredentials) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // for debug only
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if errors.Is(err, http.ErrNoCookie) {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token := cookie.Value

	err = ah.sessionService.DeleteByToken(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie = &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}
