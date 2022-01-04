package selfmerging

import (
	"swe-dashboard/internal/models"
	"testing"
	"time"
)

type mockSCM struct {
}

func (m *mockSCM) ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error) {
	patreides := &models.User{ID: 1, Username: "patreides", Name: "Paul Atreides", CreatedAt: time.Now()}
	jatreides := &models.User{ID: 2, Username: "jatreides", Name: "Lady Jessica Atreides", CreatedAt: time.Now()}
	chani := &models.User{ID: 3, Username: "chani", Name: "Chani", CreatedAt: time.Now()}
	latreides := &models.User{ID: 4, Username: "latreides", Name: "Duke Leto Atreides", CreatedAt: time.Now()}
	didaho := &models.User{ID: 5, Username: "didaho", Name: "Duncan Idaho", CreatedAt: time.Now()}

	return []models.MergeRequest{
		{ID: 1, Author: patreides, MergedBy: patreides, State: "merged"},
		{ID: 2, Author: patreides, MergedBy: patreides, State: "merged"},
		{ID: 3, Author: patreides, MergedBy: patreides, State: "merged"},
		{ID: 4, Author: jatreides, MergedBy: chani, State: "merged"},
		{ID: 5, Author: jatreides, MergedBy: latreides, State: "merged"},
		{ID: 6, Author: latreides, MergedBy: latreides, State: "merged"},
		{ID: 7, Author: didaho, MergedBy: didaho, State: "merged"},
		{ID: 8, Author: didaho, MergedBy: didaho, State: "merged"},
	}, nil
}

func TestSelfMergingUsers(t *testing.T) {
	mockscm := &mockSCM{}
	service := selfMerging{
		scm: mockscm,
	}

	selfmruserscout := 3
	list, err := service.GetSelfMergingUsers("", "", 0)
	if err != nil {
		t.Error(err)
	}

	if len(list) != selfmruserscout {
		t.Errorf("unexpeced list count. got: %d, want: %d", len(list), selfmruserscout)
	}

	user1 := list[0]
	if user1.ID != 1 {
		t.Errorf("unexpected first user. got: %d, want: 1", user1.ID)
	}

	if user1.Count != 3 {
		t.Errorf("unexpected MR count for %s. got: %f, want: 3", user1.Username, user1.Count)
	}
}
