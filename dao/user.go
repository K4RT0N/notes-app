package dao

import (
	"database/sql"
	"errors"
	"log"
	"notesapp/apperrors"
	"notesapp/entity"
	"notesapp/service"
)

type UserRepository struct {
	db Executor
}

func (ur UserRepository) Insert(userEntity entity.UserEntity) (entity.UserEntity, error) {
	row := ur.db.QueryRow("INSERT INTO users (name) VALUES ($1) RETURNING id, name", userEntity.Name)
	result := entity.UserEntity{}
	err := row.Scan(&result.Id, &result.Name)
	if err == nil {
		return result, nil
	}
	log.Println(err)
	return result, apperrors.ErrDatabaseQueryFailed
}

func (ur UserRepository) Update(userEntity entity.UserEntity) (entity.UserEntity, error) {
	row := ur.db.QueryRow("UPDATE users SET name=$1 WHERE id=$2 RETURNING id, name", userEntity.Name, userEntity.Id)
	result := entity.UserEntity{}
	err := row.Scan(&result.Id, &result.Name)
	if err == nil {
		return result, nil
	}
	log.Println(err)
	return result, apperrors.ErrDatabaseQueryFailed
}

func (ur UserRepository) FindById(userId int) (entity.UserEntity, error) {
	row := ur.db.QueryRow("SELECT * FROM users WHERE id=$1", userId)
	result := entity.UserEntity{}
	err := row.Scan(&result.Id, &result.Name)
	if err == nil {
		return result, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return entity.UserEntity{}, apperrors.ErrNoSuchUser
	}
	log.Println(err)
	return entity.UserEntity{}, apperrors.ErrDatabaseQueryFailed
}

func (ur UserRepository) FindAll() ([]entity.UserEntity, error) {
	rows, err := ur.db.Query("SELECT id, name FROM users")
	if err != nil {
		return nil, apperrors.ErrDatabaseQueryFailed
	}
	defer rows.Close()

	users := []entity.UserEntity{}
	for rows.Next() {
		user := entity.UserEntity{}
		err = rows.Scan(&user.Id, &user.Name)
		if err != nil {
			log.Println(err)
			return nil, apperrors.ErrDatabaseQueryFailed
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, apperrors.ErrDatabaseQueryFailed
	}

	return users, nil
}

func (ur UserRepository) Delete(userId int) error {
	_, err := ur.db.Exec("DELETE FROM users WHERE id=$1", userId)
	if err == nil {
		return nil
	}
	log.Println(err)
	return apperrors.ErrDatabaseQueryFailed
}

func NewUserRepository(db Executor) service.UserRepository {
	return &UserRepository{db: db}
}
