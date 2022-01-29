package models

import "time"

type Repo struct {
	ID             int
	CreatorID      int
	Name           string
	Description    string
	MRs            MergeRequests
	LastActivityAt *time.Time
	CreatedAt      *time.Time
	CommitCount    int
}
