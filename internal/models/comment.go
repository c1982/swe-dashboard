package models

import "time"

type Comment struct {
	ID           int    `json:"id"`
	Body         string `json:"body"`
	Title        string `json:"title"`
	NoteableType string
	FileName     string
	System       bool      `json:"system"`
	Resolvable   bool      `json:"resolvable"`
	Resolved     bool      `json:"resolved"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
	ApprovedNote bool
	Author       User
	ResolvedBy   User
}
