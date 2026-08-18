package entity

type AuthDataEntity struct {
	Login            string
	SaltPasswordHash string
	UserId           int
}
