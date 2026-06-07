package models

import "time"

type Wallet struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Balance   float64   `db:"balance"`
	CreatedAt time.Time `db:"created_at"`
}
