package dto

import (
	"notesapp/entity"
	"time"
)

type NoteResponse struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorId  int       `json:"author_id"`
	IsPrivate bool      `json:"is_private"`
	CreatedAt time.Time `json:"created_at"`
}

func ToNoteResponse(noteEntity entity.NoteEntity) NoteResponse {
	return NoteResponse{
		Id:        noteEntity.Id,
		Title:     noteEntity.Title,
		Content:   noteEntity.Content,
		AuthorId:  noteEntity.AuthorId,
		IsPrivate: noteEntity.IsPrivate,
		CreatedAt: noteEntity.CreatedAt,
	}
}
