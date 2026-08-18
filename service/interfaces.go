package service

import "notesapp/entity"

type TransactionManager interface {
	InTransaction(func(Transaction) error) error
}

type Transaction interface {
	User() UserRepository
	Session() SessionRepository
	AuthData() AuthDataRepository
}

type AuthDataRepository interface {
	Insert(entity.AuthDataEntity) (entity.AuthDataEntity, error)
	FindByLogin(string) (entity.AuthDataEntity, error)
}

type SessionRepository interface {
	Insert(entity.SessionEntity) (entity.SessionEntity, error)
	Update(entity.SessionEntity) (entity.SessionEntity, error)
	FindByTokenHash([]byte) (entity.SessionEntity, error)
	FindAll() ([]entity.SessionEntity, error)
	DeleteByTokenHash([]byte) error
}

type UserRepository interface {
	Insert(entity.UserEntity) (entity.UserEntity, error)
	Update(entity.UserEntity) (entity.UserEntity, error)
	FindById(int) (entity.UserEntity, error)
	FindAll() ([]entity.UserEntity, error)
	Delete(int) error
}

type NoteRepository interface {
	Insert(entity.NoteEntity) (entity.NoteEntity, error)
	FindById(int) (entity.NoteEntity, error)
	FindByAuthorId(int) ([]entity.NoteEntity, error)
	FindPublicByAuthorId(int) ([]entity.NoteEntity, error)
	DeleteById(int) error
}
