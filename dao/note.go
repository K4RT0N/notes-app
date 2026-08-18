package dao

import (
	"database/sql"
	"errors"
	"notesapp/apperrors"
	"notesapp/entity"
	"notesapp/service"
)

type NoteRepository struct {
	executor Executor
}

func NewNoteRepository(executor Executor) service.NoteRepository {
	return &NoteRepository{
		executor: executor,
	}
}

func (nr NoteRepository) Insert(note entity.NoteEntity) (entity.NoteEntity, error) {
	row := nr.executor.QueryRow("INSERT INTO notes (title, content, author_id, is_private, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, title, content, author_id, is_private, created_at",
		note.Title, note.Content, note.AuthorId, note.IsPrivate, note.CreatedAt)
	newNote := entity.NoteEntity{}
	err := row.Scan(&newNote.Id, &newNote.Title, &newNote.Content, &newNote.AuthorId, &newNote.IsPrivate, &newNote.CreatedAt)
	if err != nil {
		return entity.NoteEntity{}, apperrors.ErrDatabaseQueryFailed
	}
	return newNote, nil
}

func (nr NoteRepository) FindById(noteId int) (entity.NoteEntity, error) {
	row := nr.executor.QueryRow("SELECT id, title, content, author_id, is_private, created_at FROM notes WHERE id=$1", noteId)
	note := entity.NoteEntity{}
	err := row.Scan(&note.Id, &note.Title, &note.Content, &note.AuthorId, &note.IsPrivate, &note.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.NoteEntity{}, apperrors.ErrNoSuchNote
	}
	if err != nil {
		return entity.NoteEntity{}, apperrors.ErrDatabaseQueryFailed
	}
	return note, nil
}

func (nr NoteRepository) FindByAuthorId(authorId int) ([]entity.NoteEntity, error) {
	rows, err := nr.executor.Query("SELECT id, title, content, author_id, is_private, created_at FROM notes WHERE author_id=$1", authorId)
	if err != nil {
		return nil, apperrors.ErrDatabaseQueryFailed
	}
	defer rows.Close()
	notes := []entity.NoteEntity{}
	for rows.Next() {
		note := entity.NoteEntity{}
		err := rows.Scan(&note.Id, &note.Title, &note.Content, &note.AuthorId, &note.IsPrivate, &note.CreatedAt)
		if err != nil {
			return nil, apperrors.ErrDatabaseQueryFailed
		}
		notes = append(notes, note)
	}
	if rows.Err() != nil {
		return nil, apperrors.ErrDatabaseQueryFailed
	}
	return notes, nil
}

func (nr NoteRepository) FindPublicByAuthorId(authorId int) ([]entity.NoteEntity, error) {
	rows, err := nr.executor.Query("SELECT id, title, content, author_id, is_private, created_at FROM notes WHERE author_id=$1 AND is_private=false", authorId)
	if err != nil {
		return nil, apperrors.ErrDatabaseQueryFailed
	}
	defer rows.Close()
	notes := []entity.NoteEntity{}
	for rows.Next() {
		note := entity.NoteEntity{}
		err := rows.Scan(&note.Id, &note.Title, &note.Content, &note.AuthorId, &note.IsPrivate, &note.CreatedAt)
		if err != nil {
			return nil, apperrors.ErrDatabaseQueryFailed
		}
		notes = append(notes, note)
	}
	if rows.Err() != nil {
		return nil, apperrors.ErrDatabaseQueryFailed
	}
	return notes, nil
}

func (nr NoteRepository) DeleteById(noteId int) error {
	_, err := nr.executor.Exec("DELETE FROM notes WHERE id=$1", noteId)
	if err != nil {
		return apperrors.ErrDatabaseQueryFailed
	}
	return nil
}
