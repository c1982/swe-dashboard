package unreviewedmergerequests

import (
	"swe-dashboard/internal/models"
	"testing"
)

type mockSCM struct {
}

func (m *mockSCM) ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error) {
	return
}

func (m *mockSCM) ListMergeRequestNotes(projectID int, mergeRequestID int) (comments []*models.Comment, err error) {

	return
}

func (m *mockSCM) GetRepository(projectID int) (repository models.Repo, err error) {
	return
}

func TestIsUnreviewed(t *testing.T) {
	scm := &mockSCM{}
	srv := unreviewedMergeRequests{
		scm: scm,
	}

	checks := []struct {
		input  []*models.Comment
		expect bool
	}{
		{[]*models.Comment{{ID: 1, System: true, ApprovedNote: false}}, true},
		{[]*models.Comment{{ID: 2, System: true, ApprovedNote: true}}, false},
		{[]*models.Comment{{ID: 3, System: false, ApprovedNote: false}}, true},
	}

	for _, c := range checks {
		isunreviewed := srv.isUnreviewed(c.input)
		if isunreviewed != c.expect {
			t.Errorf("invalid unreviewed result. got: %v, want: %v", isunreviewed, c.expect)
		}
	}
}
