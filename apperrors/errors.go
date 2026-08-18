package apperrors

import "errors"

var (
	ErrNoSuchUser     = errors.New("no such user")
	ErrNoSuchAuthData = errors.New("no such login")
	ErrNoSuchSession  = errors.New("no such session")
	ErrNoSuchNote     = errors.New("no such note")
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrWrongCredentials  = errors.New("wrong credentials")
	ErrSessionExpired    = errors.New("session expired")
	ErrAccessDenied      = errors.New("access denied")
)

var (
	ErrTransactionBeginFailed = errors.New("transaction begin failed")
	ErrCryptoFuncFailed       = errors.New("crypto function failed")
	ErrDatabaseQueryFailed    = errors.New("database query failed")
)
