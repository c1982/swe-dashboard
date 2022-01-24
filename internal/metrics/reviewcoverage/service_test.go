package reviewcoverage

import (
	"swe-dashboard/internal/models"
	"testing"
)

func TestCountOfChangedFiles(t *testing.T) {
	changes := []*models.MergeRequestChanges{
		{OldPath: "a.sh", NewPath: "a.sh"},
		{OldPath: "a.sh", NewPath: "a.sh"},
		{OldPath: "b.sh", NewPath: "b.sh"},
		{OldPath: "b.sh", NewPath: "b.sh"},
		{OldPath: "c.sh", NewPath: "c.sh"},
		{OldPath: "d.sh", NewPath: "e.sh"},
	}

	srv := &reviewCoverage{}
	count := srv.countOfChangedFiles(changes)
	want := 5
	if count != want {
		t.Errorf("invalid changes count. got :%d, want: %d", count, want)
	}
}

func TestCountOfCommentedFiles(t *testing.T) {

	comments := []*models.Comment{
		{System: true},
		{System: false, FileName: ""},
		{System: false, FileName: "up.rb"},
		{System: true, FileName: "up.rb"},
	}

	srv := &reviewCoverage{}
	count := srv.countOfCommentedFiles(comments)
	want := 1
	if count != want {
		t.Errorf("invalid changes count. got :%d, want: %d", count, want)
	}
}
