package cycletime

import (
	"swe-dashboard/internal/models"
	"testing"
	"time"
)

func TestMergeRequestFirstCommit(t *testing.T) {
	commitID := "4"
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
