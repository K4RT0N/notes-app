package entity

import "time"

type NoteEntity struct {
	Id        int
	Title     string
	Content   string
	AuthorId  int
	IsPrivate bool
	CreatedAt time.Time
}
