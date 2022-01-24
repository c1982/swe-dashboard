package defectrate

import (
	"testing"
)

func TestIsDefectMergeRequest(t *testing.T) {
	tests := []struct {
		Expect bool
		Title  string
	}{
		{false, "feature/create-backend-field"},
		{true, "fix/career-unlock-control-navigator"},
		{true, "bugfix/remove-unnecessary"},
		{true, "Fix/deleting-sequence-calls"},
		{true, "refactor-deleting-sequence-calls"},
		{true, "revert-x"},
		{true, "repair-x"},
		{true, "bug-x"},
		{false, "merge-x"},
		{true, "fixfix"},
	}

	s := &defectrate{labels: defectPrefixes}
	for _, v := range tests {
		if v.Expect != s.isDefectMergeRequest(v.Title) {
			t.Errorf("title not detected. Title: %s", v.Title)
		}
	}
}
