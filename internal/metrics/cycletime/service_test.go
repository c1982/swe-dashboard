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

func TestSub(t *testing.T) {
	mergerequestopentime := time.Date(2021, 12, 13, 10, 0, 0, 0, time.UTC)
	commitopentime := time.Date(2021, 12, 13, 9, 0, 0, 0, time.UTC)
	duration := mergerequestopentime.Sub(commitopentime)
	unix := mergerequestopentime.Unix() - commitopentime.Unix()

	t.Log(duration)
	t.Log(int64(duration))
	t.Log(float64(int64(duration)) / float64(time.Nanosecond))
	t.Log(float32(unix))
}
