package models

import "time"

type UserCount struct {
	ID       int
	Name     string
	Username string
	Count    float64
	Date     time.Time
}
