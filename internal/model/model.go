package model

import "time"

type User struct {
	ID           int
	Username     string
	Email        string
	PasswordHash []byte
	RegisteredAt time.Time
	HiredAt      time.Time
}

type Transaction struct {
	ID         int
	SenderID   int
	ReceiverID int
	Amount     float64
	TypeID     int
	StatusID   int
	CreatedAt  time.Time
}
