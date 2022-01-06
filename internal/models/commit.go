package models

import "time"

type Commit struct {
	ID             string    `json:"id"`
	ShortID        string    `json:"short_id"`
	Title          string    `json:"title"`
	AuthorName     string    `json:"author_name"`
	CommitterName  string    `json:"committer_name"`
	CommitterEmail string    `json:"committer_email"`
	CommittedDate  time.Time `json:"committed_date"`
	CreatedAt      time.Time `json:"created_at"`
	Message        string    `json:"message"`
	ProjectID      int       `json:"project_id"`
}
