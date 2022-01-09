package models

import "time"

type MergeRequest struct {
	ID           int                    `json:"id"`
	IID          int                    `json:"iid"`
	TargetBranch string                 `json:"target_branch"`
	SourceBranch string                 `json:"source_branch"`
	ProjectID    int                    `json:"project_id"`
	Title        string                 `json:"title"`
	State        string                 `json:"state"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	Assignee     *User                  `json:"assignee"`
	Assignees    []*User                `json:"assignees"`
	Reviewers    []*User                `json:"reviewers"`
	Author       *User                  `json:"author"`
	MergedBy     *User                  `json:"merged_by"`
	MergedAt     *time.Time             `json:"merged_at"`
	ClosedBy     *User                  `json:"closed_by"`
	ClosedAt     *time.Time             `json:"closed_at"`
	Changes      []*MergeRequestChanges `json:"changes"`
	Draft        bool                   `json:"draft"`
}

type MergeRequestChanges struct {
	OldPath     string `json:"old_path"`
	NewPath     string `json:"new_path"`
	AMode       string `json:"a_mode"`
	BMode       string `json:"b_mode"`
	Diff        string `json:"diff"`
	NewFile     bool   `json:"new_file"`
	RenamedFile bool   `json:"renamed_file"`
	DeletedFile bool   `json:"deleted_file"`
}
