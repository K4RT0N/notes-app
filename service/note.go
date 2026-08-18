package service

import (
	"crypto/sha256"
	"notesapp/apperrors"
	"notesapp/entity"
	"time"
)

type NoteService struct {
	noteRepository    NoteRepository
	sessionRepository SessionRepository
}

func NewNoteService(noteRepository NoteRepository, sessionRepository SessionRepository) *NoteService {
	return &NoteService{
		noteRepository:    noteRepository,
		sessionRepository: sessionRepository,
	}
}

func (ns NoteService) GetById(token string, noteId int) (entity.NoteEntity, error) {
	note, err := ns.noteRepository.FindById(noteId)
	if err != nil {
		return entity.NoteEntity{}, err
	}

	if !note.IsPrivate {
		return note, nil
	}

	tokenHash := sha256.Sum256([]byte(token))
	session, err := ns.sessionRepository.FindByTokenHash(tokenHash[:])
	if err != nil {
		return entity.NoteEntity{}, err
	}
	if session.ExpiresAt.Before(time.Now()) {
		return entity.NoteEntity{}, apperrors.ErrSessionExpired
	}

	if session.UserId == note.AuthorId {
		return note, nil
	}

	return entity.NoteEntity{}, apperrors.ErrAccessDenied
}

func (ns NoteService) CreateNote(token string, noteEntity entity.NoteEntity) (entity.NoteEntity, error) {
	tokenHash := sha256.Sum256([]byte(token))
	session, err := ns.sessionRepository.FindByTokenHash(tokenHash[:])
	if err != nil {
		return entity.NoteEntity{}, err
	}
	if session.ExpiresAt.Before(time.Now()) {
		return entity.NoteEntity{}, apperrors.ErrSessionExpired
	}
	noteEntity.AuthorId = session.UserId
	newNote, err := ns.noteRepository.Insert(noteEntity)
	if err != nil {
		return entity.NoteEntity{}, err
	}
	return newNote, nil
}

func (ns NoteService) GetMyNotes(token string) ([]entity.NoteEntity, error) {
	tokenHash := sha256.Sum256([]byte(token))
	session, err := ns.sessionRepository.FindByTokenHash(tokenHash[:])
	if err != nil {
		return nil, err
	}
	if session.ExpiresAt.Before(time.Now()) {
		return nil, apperrors.ErrSessionExpired
	}
	userId := session.UserId
	userNotes, err := ns.noteRepository.FindByAuthorId(userId)
	if err != nil {
		return nil, err
	}
	return userNotes, nil
}

func (ns NoteService) GetUserNotes(userId int) ([]entity.NoteEntity, error) {
	publicNotes, err := ns.noteRepository.FindPublicByAuthorId(userId)
	if err != nil {
		return nil, err
	}

	return publicNotes, nil
}

func (ns NoteService) DeleteById(token string, noteId int) error {
	tokenHash := sha256.Sum256([]byte(token))
	session, err := ns.sessionRepository.FindByTokenHash(tokenHash[:])
	if err != nil {
		return err
	}
	if session.ExpiresAt.Before(time.Now()) {
		return apperrors.ErrSessionExpired
	}
	note, err := ns.noteRepository.FindById(noteId)
	if err != nil {
		return err
	}
	if session.UserId != note.AuthorId {
		return apperrors.ErrAccessDenied
	}
	return ns.noteRepository.DeleteById(noteId)
}
