package models

import "time"

type Repo struct {
	ID             int
	Name           string
	Description    string
	LastActivityAt *time.Time
	CreatorID      int
	MRs            MergeRequests
}
