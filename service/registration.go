package service

import (
	"errors"
	"notesapp/apperrors"
	"notesapp/entity"

	"github.com/alexedwards/argon2id"
)

type RegistrationService struct {
	tm                 TransactionManager
	authDataRepository AuthDataRepository
}

func NewRegistrationService(tm TransactionManager, authDataDAO AuthDataRepository) *RegistrationService {
	return &RegistrationService{
		tm:                 tm,
		authDataRepository: authDataDAO,
	}
}

func (rs *RegistrationService) Registration(login, password string) error {
	_, err := rs.authDataRepository.FindByLogin(login)
	if err == nil {
		return apperrors.ErrUserAlreadyExists
	}
	if !errors.Is(err, apperrors.ErrNoSuchAuthData) {
		return err
	}

	user := entity.UserEntity{
		Name: login,
	}
	saltPasswordHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return apperrors.ErrCryptoFuncFailed
	}
	authData := entity.AuthDataEntity{
		Login:            login,
		SaltPasswordHash: saltPasswordHash,
	}

	err = rs.tm.InTransaction(func(transaction Transaction) error {
		user, err := transaction.User().Insert(user)
		if err != nil {
			return err
		}

		authData.UserId = user.Id
		_, err = transaction.AuthData().Insert(authData)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
