package models

import "time"

type Wallet struct {
	ID             int       `db:"id"`
	UserID         int       `db:"user_id"`
	InitialBalance float64   `db:"initial_balance"`
	CreatedAt      time.Time `db:"created_at"`
}
