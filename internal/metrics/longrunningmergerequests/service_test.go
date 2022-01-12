package longrunningmergerequests

import (
	"swe-dashboard/internal/models"
	"testing"
	"time"
)

func TestMergeRequestLastActivity(t *testing.T) {
	srv := longRunningMergerequests{}
	comments := []*models.Comment{
		{ID: 1, System: true, CreatedAt: time.Date(2021, 12, 1, 0, 0, 0, 0, time.Local), ApprovedNote: false},
		{ID: 2, System: false, CreatedAt: time.Date(2021, 13, 1, 0, 0, 0, 0, time.Local), ApprovedNote: false},
		{ID: 3, System: false, CreatedAt: time.Date(2021, 14, 1, 0, 0, 0, 0, time.Local), ApprovedNote: false},
		{ID: 4, System: true, CreatedAt: time.Date(2021, 15, 1, 0, 0, 0, 0, time.Local), ApprovedNote: true},
	}

	commentID := 4
	comment := srv.mergeRequestLastActivity(comments)
	if comment.ID != commentID {
		t.Errorf("invalid comment ID. got: %d, want: %d", comment.ID, commentID)
	}
}
