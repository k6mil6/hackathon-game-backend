package model

type User struct {
	ID           int
	Username     string
	PasswordHash []byte
	CreatedAt    string
}
