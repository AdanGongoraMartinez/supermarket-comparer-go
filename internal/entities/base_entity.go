package entities

import "time"

type BaseEntity struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

