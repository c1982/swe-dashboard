package models

import "time"

type Comment struct {
	ID         int        `json:"id"`
	Body       string     `json:"body"`
	Title      string     `json:"title"`
	System     bool       `json:"system"`
	Resolvable bool       `json:"resolvable"`
	Resolved   bool       `json:"resolved"`
	ExpiresAt  *time.Time `json:"expires_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	CreatedAt  *time.Time `json:"created_at"`
	Author     User
	ResolvedBy User
}
