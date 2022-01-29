package models

import "time"

type Commit struct {
	ID             string    `json:"id"`
	ShortID        string    `json:"short_id"`
	Title          string    `json:"title"`
	AuthorName     string    `json:"author_name"`
	AuthorEmail    string    `json:"author_email"`
	CommitterName  string    `json:"committer_name"`
	CommitterEmail string    `json:"committer_email"`
	Message        string    `json:"message"`
	Additions      int       `json:"additions"`
	Deletions      int       `json:"deletions"`
	Total          int       `json:"total"`
	ProjectID      int       `json:"project_id"`
	CommittedDate  time.Time `json:"committed_date"`
	CreatedAt      time.Time `json:"created_at"`
}
