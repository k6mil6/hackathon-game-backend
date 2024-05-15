package model

import "time"

type User struct {
	ID           int
	Username     string
	Balance      float64
	Email        string
	PasswordHash []byte
	ClassID      int
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

type ShopItem struct {
	ID          int
	Name        string
	Description string
	Price       float64
	InStock     int
}

type Purchase struct {
	ID         int
	ShopItemID int
	BuyerID    int
	CreatedAt  time.Time
}

type Admin struct {
	ID           int
	Username     string
	Email        string
	PasswordHash []byte
	RegisteredAt time.Time
	RegisteredBy int
	RoleID       int
}

type Task struct {
	ID         int
	Name       string
	StatusID   int
	Amount     float64
	CreatedAt  time.Time
	CreatedBy  int
	ForGroupID int
	UserID     int
}

type Business struct {
	ID      int
	Name    string
	TypeID  int
	OwnerID int
	Price   float64
}

type BusinessType struct {
	ID          int
	Name        string
	Description string
	Profit      float64
}
