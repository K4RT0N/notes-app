package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"notesapp/apperrors"
	"notesapp/entity"
	"time"

	"github.com/alexedwards/argon2id"
)

type SessionService struct {
	sessionDAO  SessionRepository
	authDataDAO AuthDataRepository
}

func NewSessionService(sessionDAO SessionRepository, authDataDAO AuthDataRepository) *SessionService {
	return &SessionService{
		sessionDAO:  sessionDAO,
		authDataDAO: authDataDAO,
	}
}

func (ss SessionService) InitializeSession(login, password string) (string, error) {
	authData, err := ss.authDataDAO.FindByLogin(login)
	if errors.Is(err, apperrors.ErrNoSuchAuthData) {
		return "", apperrors.ErrWrongCredentials
	}

	res, err := argon2id.ComparePasswordAndHash(password, authData.SaltPasswordHash)
	if err != nil {
		return "", apperrors.ErrCryptoFuncFailed
	}
	if !res {
		return "", apperrors.ErrWrongCredentials
	}

	tokenRaw := make([]byte, 16)
	_, err = rand.Read(tokenRaw)
	if err != nil {
		return "", apperrors.ErrCryptoFuncFailed
	}
	token := hex.EncodeToString(tokenRaw)
	tokenHash := sha256.Sum256([]byte(token))
	now := time.Now()
	session := entity.SessionEntity{
		TokenHash: tokenHash[:],
		CreatedAt: now,
		ExpiresAt: now.AddDate(0, 0, 7),
		UserId:    authData.UserId,
	}
	_, err = ss.sessionDAO.Insert(session)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (ss SessionService) DeleteByToken(token string) error {
	tokenHash := sha256.Sum256([]byte(token))
	return ss.sessionDAO.DeleteByTokenHash(tokenHash[:])
}
