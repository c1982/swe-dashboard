package assetworkingtime

import (
	"swe-dashboard/internal/models"
	"testing"
	"time"
)

type MockSCM struct{}

func (m *MockSCM) ListProjects() (repositories []*models.Repo, err error) {
	return []*models.Repo{
		{ID: 26, CreatorID: 1, Name: "project-1", Description: "",
			MRs:            []models.MergeRequest{},
			LastActivityAt: m.Time("2022-10-01 13:00:00"),
			CreatedAt:      m.Time("2022-10-01 13:00:00"),
			CommitCount:    1},
	}, nil
}

func (m *MockSCM) ListCommits(projectID int, createdafterday int) (commits []*models.Commit, err error) {
	return []*models.Commit{
		{ID: "1", ShortID: "1", Title: "", AuthorName: "", AuthorEmail: "", CommitterName: "",
			CommitterEmail: "", Message: "", Additions: 0, Deletions: 0, Total: 0, ProjectID: 26,
			CommittedDate: m.Time2("2022-10-05 13:00"), CreatedAt: m.Time2("2022-10-05 13:00")},
		{ID: "1", ShortID: "2", Title: "", AuthorName: "", AuthorEmail: "", CommitterName: "",
			CommitterEmail: "", Message: "", Additions: 0, Deletions: 0, Total: 0, ProjectID: 26,
			CommittedDate: m.Time2("2022-10-05 15:00"), CreatedAt: m.Time2("2022-10-05 15:00")},
		{ID: "1", ShortID: "3", Title: "", AuthorName: "", AuthorEmail: "", CommitterName: "",
			CommitterEmail: "", Message: "", Additions: 0, Deletions: 0, Total: 0, ProjectID: 26,
			CommittedDate: m.Time2("2022-10-05 18:00"), CreatedAt: m.Time2("2022-10-05 18:00")},
	}, nil
}

func (m *MockSCM) CommitChanges(projectID int, commitID string) (changes []*models.Change, err error) {
	list := map[string][]*models.Change{
		"1": {
			{ProjectID: 26, Name: "a.png", Weight: 10},
			{ProjectID: 26, Name: "a.png", Weight: 10},
		},
		"2": {
			{ProjectID: 26, Name: "a.png", Weight: 10},
			{ProjectID: 26, Name: "b.psd", Weight: 30},
		},
		"3": {
			{ProjectID: 26, Name: "b.psd", Weight: 30},
		},
	}
	return list[commitID], nil
}

func (m *MockSCM) Time(v string) *time.Time {
	t, _ := time.Parse(time.RFC3339, v)
	p := new(time.Time)
	*p = t
	return p
}

func (m *MockSCM) Time2(v string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04", v)
	return t
}

func TestCalculations(t *testing.T) {
	scm := &MockSCM{}
	svc := NewAssetWorkingTimeService(scm, ".psd", ".png")
	err := svc.CalculateChanges()
	if err != nil {
		t.Error(err)
	}

	t.Run("Weights", func(r *testing.T) {
		weights := svc.Weights()
		if len(weights) != 2 {
			t.Errorf("unexpected weight. got: %d, want: 2", len(weights))
		}

		if weights[0].Name1 != "a.png" {
			t.Errorf("index name error, got: %s, want: a.png", weights[0].Name1)
		}

		if weights[0].Count != 30 {
			t.Errorf("asset weight error. got: %f, want: 30", weights[0].Count)
		}

		if weights[1].Name1 != "b.psd" {
			t.Errorf("index name error, got: %s, want: a.png", weights[1].Name1)
		}

		if weights[1].Count != 60 {
			t.Errorf("asset weight error. got: %f, want: 30", weights[1].Count)
		}
	})

	t.Run("Iterations", func(r *testing.T) {
		iterations := svc.Iterations()
		if len(iterations) != 2 {
			t.Errorf("unexpected iterations. got: %d, want: 2", len(iterations))
		}

		if iterations[0].Count != 3 {
			t.Errorf("iteration error. got: %f, want: 3", iterations[0].Count)
		}
		if iterations[1].Count != 2 {
			t.Errorf("iteration error. got: %f, want: 2", iterations[1].Count)
		}
	})

	t.Run("WorkingHours", func(r *testing.T) {
		hours := svc.WorkingHours()
		if len(hours) != 2 {
			t.Errorf("unexpected hours. got: %d, want: 2", len(hours))
		}
		if hours[0].Count != 2 {
			t.Errorf("working hours error. got: %f, want: 2", hours[0].Count)
		}
		if hours[1].Count != 3 {
			t.Errorf("working hours error. got: %f, want: 3", hours[1].Count)
		}
	})
}
