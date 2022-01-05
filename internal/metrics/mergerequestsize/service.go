package mergerequestsize

import (
	"strings"
	"swe-dashboard/internal/models"
	"time"
)

type SCM interface {
	GetRepository(projectID int) (repository models.Repo, err error)
	GetMergeRequestChanges(projectID int, mergeRequestID int) (mergerequest models.MergeRequest, err error)
	ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error)
}

type MergeRequestSizeService interface {
	MergeRequestSizes() (sizes []models.ItemCount, err error)
}

func NewMergeRequestSizeService(scm SCM) MergeRequestSizeService {
	m := &mergeRequestSizes{
		scm: scm,
	}
	return m
}

type mergeRequestSizes struct {
	scm           SCM
	mergerequests models.MergeRequests
}

func (m *mergeRequestSizes) MergeRequestSizes() (sizes []models.ItemCount, err error) {
	m.mergerequests, err = m.scm.ListMergeRequest("merged", "all", time.Now().Day())
	if err != nil {
		return sizes, err
	}

	sizes = []models.ItemCount{}
	repositories := m.mergerequests.GroupByRepositories()
	for i := 0; i < len(repositories); i++ {
		repo, err := m.scm.GetRepository(repositories[i].ID)
		if err != nil {
			return sizes, err
		}
		repo.MRs = repositories[i].MRs

		for n := 0; n < len(repo.MRs); n++ {
			mr := repo.MRs[n]
			singlemr, err := m.scm.GetMergeRequestChanges(repo.ID, mr.IID)
			if err != nil {
				return sizes, err
			}

			newline, deletedline := m.calculateChanges(singlemr.Changes)
			sizes = append(sizes, models.ItemCount{
				Name:  repo.Name,
				Name1: singlemr.Title,
				Count: newline + deletedline,
			})
		}
	}

	return sizes, err
}

func (m *mergeRequestSizes) calculateChanges(changes []*models.MergeRequestChanges) (newline, deletedline float64) {
	for c := 0; c < len(changes); c++ {
		change := changes[c]
		lines := strings.Split(change.Diff, "\n")
		for l := 0; l < len(lines); l++ {
			line := strings.TrimPrefix(lines[l], "\"")
			if strings.HasPrefix(line, "+") {
				newline = newline + 1
			}

			if strings.HasPrefix(line, "-") {
				deletedline = deletedline + 1
			}
		}
	}

	return newline, deletedline
}
