package model

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Id     uuid.UUID `json:"id"`
	UserId int64     `json:"user_id"`

	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`

	CreatedAt  time.Time `json:"created_at"`
	LastSeenAt time.Time `json:"last_seen_at"`
}
