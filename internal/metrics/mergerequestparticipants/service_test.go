package mergerequestparticipants

import (
	"fmt"
	"swe-dashboard/internal/models"
	"testing"
)

func TestCountParticipants(t *testing.T) {
	src := mergeRequestParticipants{}
	counts := map[int]*models.ItemCount{}
	participants := []*models.User{
		{ID: 1, Username: "patreides", Name: "Paul Atreides"},
		{ID: 1, Username: "patreides", Name: "Paul Atreides"},
		{ID: 2, Username: "jatreides", Name: "Lady Jessica Atreides"},
		{ID: 3, Username: "chani", Name: "Chani"},
		{ID: 4, Username: "latreides", Name: "Duke Leto Atreides"},
		{ID: 5, Username: "didaho", Name: "Duncan Idaho"},
		{ID: 5, Username: "didaho", Name: "Duncan Idaho"},
		{ID: 5, Username: "didaho", Name: "Duncan Idaho"},
	}
	author := &models.User{ID: 1, Username: "patreides", Name: "Paul Atreides"}

	src.countParticipants(author, counts, participants)

	if len(counts) != 4 {
		t.Errorf("invalid counts. got: %d, want: 4", len(counts))
	}

	v := counts[5]
	if v.Name1 != "didaho" {
		t.Errorf("invalid user. got: %s, want: didaho", v.Name1)
	}

	for _, v := range counts {
		if v.Name != "" {
			t.Errorf("repo name must be empty")
		}
	}
}

func TestEngagementParticipantCounts(t *testing.T) {
	counts := map[int]map[int]*engageItem{}
	participants := []*models.User{
		{ID: 1, Username: "patreides", Name: "Paul Atreides"},
		{ID: 2, Username: "jatreides", Name: "Lady Jessica Atreides"},
		{ID: 3, Username: "chani", Name: "Chani"},
		{ID: 4, Username: "latreides", Name: "Duke Leto Atreides"},
		{ID: 5, Username: "didaho", Name: "Duncan Idaho"},
	}
	authors := []*models.User{
		{ID: 1, Username: "patreides", Name: "Paul Atreides"},
		{ID: 2, Username: "jatreides", Name: "Lady Jessica Atreides"},
	}

	svc := mergeRequestParticipants{}

	for i := 0; i < len(authors); i++ {
		author := authors[i]
		svc.engagementParticipantCounts(author, counts, participants)
	}

	if len(counts) != 2 {
		t.Errorf("invalid counts. got: %d, want: 2", len(counts))
	}

	a1, ok := counts[1]
	if !ok {
		t.Error("author id not found. want: 1")
	}

	a2, ok := counts[2]
	if !ok {
		t.Error("author id not found. want: 2")
	}

	if len(a1) != 4 {
		t.Errorf("invalid participant count a1. got: %d, want: 4", len(a1))
	}

	if len(a2) != 4 {
		t.Errorf("invalid participant count a2. got: %d, want: 4", len(a2))
	}
}

func TestEngagementCounts(t *testing.T) {
	counts := map[string]*engageItem{}
	mrs := []models.MergeRequest{
		{ID: 1,
			Author:   &models.User{ID: 1, Username: "patreides", Name: "Paul Atreides"},
			MergedBy: &models.User{ID: 2, Username: "jatreides", Name: "Lady Jessica Atreides"}},
		{ID: 2,
			Author:   &models.User{ID: 1, Username: "patreides", Name: "Paul Atreides"},
			MergedBy: &models.User{ID: 3, Username: "chani", Name: "Chani"}},
		{ID: 3,
			Author:   &models.User{ID: 5, Username: "didaho", Name: "Duncan Idaho"},
			MergedBy: &models.User{ID: 5, Username: "didaho", Name: "Duncan Idaho"}},
	}

	svc := mergeRequestParticipants{}
	for i := 0; i < len(mrs); i++ {
		svc.engagementCounts(mrs[i], counts)
	}

	if len(counts) != 3 {
		t.Errorf("invalid counts. got: %d, want: 3", len(counts))
	}

	v, ok := counts[fmt.Sprintf("%d-%d", 1, 2)]
	if !ok {
		t.Error("pair not found. want: 1-2")
	}

	if v.Author.ID != 1 {
		t.Errorf("invalid author id. got: %d want: 1", v.Author.ID)
	}
}
