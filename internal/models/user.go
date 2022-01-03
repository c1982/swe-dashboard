package models

import "time"

type User struct {
	ID             int
	IsAdmin        bool
	Username       string
	Name           string
	Email          string
	State          string
	AvatarURL      string
	CreatedAt      time.Time
	LastSignInAt   time.Time
	LastActivityOn time.Time
}
