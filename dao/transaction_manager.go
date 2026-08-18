package dao

import (
	"context"
	"database/sql"
	"notesapp/apperrors"
	"notesapp/service"
)

type TransactionManager struct {
	db *sql.DB
}

type Transaction struct {
	db *sql.DB
	tx *sql.Tx
}

func (transaction *Transaction) User() service.UserRepository {
	return NewUserRepository(transaction.tx)
}

func (transaction *Transaction) Session() service.SessionRepository {
	return NewSessionRepository(transaction.tx)
}

func (transaction *Transaction) AuthData() service.AuthDataRepository {
	return NewAuthDataRepository(transaction.tx)
}

func newTransaction(db *sql.DB) (*Transaction, error) {
	tx, err := db.BeginTx(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	return &Transaction{
		tx: tx,
	}, nil
}

func NewTransactionManager(db *sql.DB) service.TransactionManager {
	return &TransactionManager{
		db: db,
	}
}

func (tm *TransactionManager) InTransaction(fn func(transaction service.Transaction) error) error {
	t, err := newTransaction(tm.db)
	if err != nil {
		return apperrors.ErrTransactionBeginFailed
	}

	err = fn(t)
	if err != nil {
		t.tx.Rollback()
		return err
	}
	return t.tx.Commit()
}
