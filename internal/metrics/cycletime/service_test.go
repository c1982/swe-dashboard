package cycletime

import (
	"swe-dashboard/internal/models"
	"testing"
	"time"
)

func TestMergeRequestFirstCommit(t *testing.T) {
	commitID := "1"
	s := cycleTime{}
	commits := []*models.Commit{
		{ID: "1", CreatedAt: time.Date(2021, 12, 1, 0, 0, 0, 0, time.Local)},
		{ID: "2", CreatedAt: time.Date(2021, 12, 10, 0, 0, 0, 0, time.Local)},
		{ID: "3", CreatedAt: time.Date(2021, 12, 20, 0, 0, 0, 0, time.Local)},
		{ID: "4", CreatedAt: time.Date(2022, 1, 5, 0, 0, 0, 0, time.Local)},
	}
	commit := s.mergeRequestFirstCommit(commits)
	if commit.ID != commitID {
		t.Errorf("invalid commit ID. got: %s, want: %s", commit.ID, commitID)
	}
}

func TestMergeRequestFirstComment(t *testing.T) {
	commentID := 2
	s := cycleTime{}
	comments := []*models.Comment{
		{ID: 1, System: true, CreatedAt: time.Date(2021, 12, 1, 0, 0, 0, 0, time.Local), ApprovedNote: false},
		{ID: 2, System: false, CreatedAt: time.Date(2021, 13, 1, 0, 0, 0, 0, time.Local), ApprovedNote: false},
		{ID: 3, System: false, CreatedAt: time.Date(2021, 14, 1, 0, 0, 0, 0, time.Local), ApprovedNote: false},
		{ID: 4, System: true, CreatedAt: time.Date(2021, 15, 1, 0, 0, 0, 0, time.Local), ApprovedNote: true},
	}

	comment := s.mergeRequestFirstComment(comments)
	if comment.ID != commentID {
		t.Errorf("invalid comment ID. got: %d, want: %d", comment.ID, commentID)
	}
}

func TestMergeRequestApprovalComment(t *testing.T) {
	s := cycleTime{}
	commentID := 2
	comments := []*models.Comment{
		{ID: 1, System: true, CreatedAt: time.Date(2021, 12, 1, 0, 0, 0, 0, time.Local), ApprovedNote: false},
		{ID: 2, System: true, CreatedAt: time.Date(2021, 13, 1, 0, 0, 0, 0, time.Local), ApprovedNote: true},
		{ID: 3, System: true, CreatedAt: time.Date(2021, 14, 1, 0, 0, 0, 0, time.Local), ApprovedNote: true},
		{ID: 4, System: true, CreatedAt: time.Date(2021, 15, 1, 0, 0, 0, 0, time.Local), ApprovedNote: false},
	}
	comment := s.mergeRequestApprovalComment(comments)
	if comment.ID != commentID {
		t.Errorf("invalid comment ID. got: %d, want: %d", comment.ID, commentID)
	}
}

func TestCleanTitle(t *testing.T) {
	s := cycleTime{}
	title := "\"this is \"title\" payload {$1}"
	cleaned := s.cleanTitle(title)
	wanted := "this is title payload $1"
	if cleaned != wanted {
		t.Errorf("string does not cleared: got: %s, want: %s", cleaned, wanted)
	}
}
