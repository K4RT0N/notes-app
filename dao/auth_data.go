package dao

import (
	"database/sql"
	"log"
	"notesapp/apperrors"
	"notesapp/entity"
	"notesapp/service"
)

type AuthDataRepository struct {
	exec Executor
}

func (adr AuthDataRepository) Insert(authData entity.AuthDataEntity) (entity.AuthDataEntity, error) {
	row := adr.exec.QueryRow("INSERT INTO auth_data (login, password_hash, user_id) VALUES ($1, $2, $3) RETURNING login, password_hash, user_id",
		authData.Login, authData.SaltPasswordHash, authData.UserId)
	err := row.Scan(&authData.Login, &authData.SaltPasswordHash, &authData.UserId)
	if err == sql.ErrNoRows {
		return entity.AuthDataEntity{}, apperrors.ErrNoSuchAuthData
	}
	if err != nil {
		log.Println(err)
		return entity.AuthDataEntity{}, apperrors.ErrDatabaseQueryFailed
	}
	return authData, nil
}

func (add AuthDataRepository) FindByLogin(login string) (entity.AuthDataEntity, error) {
	row := add.exec.QueryRow("SELECT login, password_hash, user_id FROM auth_data WHERE login=$1",
		login)
	authData := entity.AuthDataEntity{}
	err := row.Scan(&authData.Login, &authData.SaltPasswordHash, &authData.UserId)
	if err == sql.ErrNoRows {
		return entity.AuthDataEntity{}, apperrors.ErrNoSuchAuthData
	}
	if err != nil {
		log.Println(err)
		return entity.AuthDataEntity{}, apperrors.ErrDatabaseQueryFailed
	}
	return authData, nil
}

func NewAuthDataRepository(exec Executor) service.AuthDataRepository {
	return &AuthDataRepository{
		exec: exec,
	}
}
