package models

import "time"

type Wallet struct {
	ID        string    `json:"id" db:"id"`
	Balance   float64   `json:"balance" db:"balance"`
	UserID    string    `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
