package model

import "time"

type User struct {
	Id int64 `json:"id"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	Email            string    `json:"email"`
	EmailConfirmed   bool      `json:"email_confirmed"`
	EmailConfirmedAt time.Time `json:"email_confirmed_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
