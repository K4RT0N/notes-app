package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"notesapp/apperrors"
	"notesapp/dto"
	"notesapp/entity"
	"notesapp/service"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type NoteHandler struct {
	noteService *service.NoteService
}

func NewNoteHandler(noteService *service.NoteService) *NoteHandler {
	return &NoteHandler{
		noteService: noteService,
	}
}

func (nh NoteHandler) GetUserNotes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIdStr := vars["userId"]
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	publicUserNotes, err := nh.noteService.GetUserNotes(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	publicUserNotesResponse := []dto.NoteResponse{}
	for _, noteEntity := range publicUserNotes {
		publicUserNotesResponse = append(publicUserNotesResponse, dto.ToNoteResponse(noteEntity))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(publicUserNotesResponse)
}

func (nh NoteHandler) GetMyNotes(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := cookie.Value
	userNotes, err := nh.noteService.GetMyNotes(token)
	if errors.Is(err, apperrors.ErrNoSuchSession) || errors.Is(err, apperrors.ErrSessionExpired) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userNotesResponse := []dto.NoteResponse{}
	for _, noteEntity := range userNotes {
		noteResponse := dto.ToNoteResponse(noteEntity)
		userNotesResponse = append(userNotesResponse, noteResponse)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userNotesResponse)
}

func (nh NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")
	isPrivateStr := r.FormValue("is_private")
	var isPrivate bool
	if isPrivateStr == "on" {
		isPrivate = true
	} else {
		isPrivate = false
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := cookie.Value
	noteEntity := entity.NoteEntity{
		Title:     title,
		Content:   content,
		IsPrivate: isPrivate,
		CreatedAt: time.Now(),
	}

	noteEntity, err = nh.noteService.CreateNote(token, noteEntity)
	if errors.Is(err, apperrors.ErrNoSuchSession) || errors.Is(err, apperrors.ErrSessionExpired) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if errors.Is(err, apperrors.ErrDatabaseQueryFailed) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (nh NoteHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	noteIdStr := vars["noteId"]
	noteId, err := strconv.Atoi(noteIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := cookie.Value

	err = nh.noteService.DeleteById(token, noteId)
	if errors.Is(err, apperrors.ErrAccessDenied) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if errors.Is(err, apperrors.ErrNoSuchSession) || errors.Is(err, apperrors.ErrSessionExpired) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if errors.Is(err, apperrors.ErrDatabaseQueryFailed) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if errors.Is(err, apperrors.ErrNoSuchNote) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (hn NoteHandler) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	noteIdStr := vars["noteId"]
	noteId, err := strconv.Atoi(noteIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := cookie.Value

	note, err := hn.noteService.GetById(token, noteId)
	if errors.Is(err, apperrors.ErrSessionExpired) || errors.Is(err, apperrors.ErrNoSuchSession) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if errors.Is(err, apperrors.ErrAccessDenied) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if errors.Is(err, apperrors.ErrNoSuchNote) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	noteResponse := dto.ToNoteResponse(note)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(noteResponse)
}
