package dao

import (
	"database/sql"
	"log"
	"notesapp/apperrors"
	"notesapp/entity"
	"notesapp/service"
)

type SessionRepository struct {
	db Executor
}

func NewSessionRepository(db Executor) service.SessionRepository {
	return &SessionRepository{db: db}
}

func (sr SessionRepository) Insert(sessionEntity entity.SessionEntity) (entity.SessionEntity, error) {
	row := sr.db.QueryRow("INSERT INTO sessions (token_hash, created_at, expires_at, user_id) VALUES ($1, $2, $3, $4) RETURNING token_hash, created_at, expires_at, user_id",
		sessionEntity.TokenHash, sessionEntity.CreatedAt, sessionEntity.ExpiresAt, sessionEntity.UserId)
	session := entity.SessionEntity{}
	err := row.Scan(&session.TokenHash, &session.CreatedAt, &session.ExpiresAt, &session.UserId)
	if err != nil {
		log.Println(err)
		return entity.SessionEntity{}, apperrors.ErrDatabaseQueryFailed
	}
	return session, nil
}

func (sr *SessionRepository) Update(sessionEntity entity.SessionEntity) (entity.SessionEntity, error) {
	row := sr.db.QueryRow("UPDATE sessions SET created_at=$1, expires_at=$2, user_id=$3 WHERE token_hash=$4 RETURNING token_hash, created_at, expires_at, user_id",
		sessionEntity.CreatedAt, sessionEntity.ExpiresAt, sessionEntity.UserId, sessionEntity.TokenHash)
	session := entity.SessionEntity{}
	err := row.Scan(&session.TokenHash, &session.CreatedAt, &session.ExpiresAt, &session.UserId)
	if err == sql.ErrNoRows {
		return entity.SessionEntity{}, apperrors.ErrNoSuchSession
	}
	if err != nil {
		log.Println(err)
		return entity.SessionEntity{}, apperrors.ErrDatabaseQueryFailed
	}
	return session, nil
}

func (sr SessionRepository) FindByTokenHash(tokenHash []byte) (entity.SessionEntity, error) {
	row := sr.db.QueryRow("SELECT token_hash, created_at, expires_at, user_id FROM sessions WHERE token_hash=$1", tokenHash)
	session := entity.SessionEntity{}
	err := row.Scan(&session.TokenHash, &session.CreatedAt, &session.ExpiresAt, &session.UserId)
	if err == sql.ErrNoRows {
		return entity.SessionEntity{}, apperrors.ErrNoSuchSession
	}
	if err != nil {
		log.Println(err)
		return entity.SessionEntity{}, apperrors.ErrDatabaseQueryFailed
	}
	return session, nil
}

func (sr SessionRepository) FindAll() ([]entity.SessionEntity, error) {
	rows, err := sr.db.Query("SELECT token_hash, created_at, expires_at, user_id FROM sessions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := []entity.SessionEntity{}
	for rows.Next() {
		session := entity.SessionEntity{}
		err := rows.Scan(&session.TokenHash, &session.CreatedAt, &session.ExpiresAt, &session.UserId)
		if err != nil {
			log.Println(err)
			return nil, apperrors.ErrDatabaseQueryFailed
		}
		sessions = append(sessions, session)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, apperrors.ErrDatabaseQueryFailed
	}
	return sessions, nil
}

func (sr SessionRepository) DeleteByTokenHash(tokenHash []byte) error {
	_, err := sr.db.Exec("DELETE FROM sessions WHERE token_hash=$1", tokenHash)
	if err != nil {
		log.Println(err)
		return apperrors.ErrDatabaseQueryFailed
	}
	return nil
}
