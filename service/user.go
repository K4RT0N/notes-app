package service

import (
	"crypto/sha256"
	"log"
	"notesapp/apperrors"
	"notesapp/entity"
	"time"
)

type UserService struct {
	userDAO    UserRepository
	sessionDAO SessionRepository
}

func NewUserService(userDAO UserRepository, sessionDAO SessionRepository) *UserService {
	return &UserService{
		userDAO:    userDAO,
		sessionDAO: sessionDAO,
	}
}

func (us *UserService) GetById(userId int) (entity.UserEntity, error) {
	return us.userDAO.FindById(userId)
}

func (us *UserService) GetMe(sessionToken string) (entity.UserEntity, error) {
	tokenHash := sha256.Sum256([]byte(sessionToken))
	session, err := us.sessionDAO.FindByTokenHash(tokenHash[:])
	if err != nil {
		return entity.UserEntity{}, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		err := us.sessionDAO.DeleteByTokenHash(tokenHash[:])
		if err != nil {
			log.Println(err)
		}
		return entity.UserEntity{}, apperrors.ErrSessionExpired
	}

	user, err := us.userDAO.FindById(session.UserId)
	if err != nil {
		return entity.UserEntity{}, err
	}

	return user, nil
}

func (us *UserService) GetAll() ([]entity.UserEntity, error) {
	return us.userDAO.FindAll()
}
