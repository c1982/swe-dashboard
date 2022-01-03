package mergerequestsize

import (
	"sort"
	"strings"
	"swe-dashboard/internal/models"
	"time"
)

type SCM interface {
	GetMergeRequestChanges(projectID int, mergeRequestID int) (mergerequest models.MergeRequest, err error)
	ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error)
}

type MergeRequestSizeService interface {
	MergeRequestSizes(state, scope string, createdafterday int) (sizes []models.ItemCount, err error)
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

func (m *mergeRequestSizes) MergeRequestSizes(state, scope string, createdafterday int) (sizes []models.ItemCount, err error) {
	m.mergerequests, err = m.scm.ListMergeRequest(state, scope, createdafterday)
	if err != nil {
		return sizes, err
	}

	sizes = []models.ItemCount{}
	repositorygroups := m.mergerequests.GroupByRepositories()
	for i := 0; i < len(repositorygroups); i++ {
		repo := repositorygroups[i]
		for n := 0; n < len(repo.MRs); n++ {
			mr := repo.MRs[n]
			singlemr, err := m.scm.GetMergeRequestChanges(repo.ID, mr.IID)
			if err != nil {
				return sizes, err
			}

			size := m.calculateChanges(mr.CreatedAt, singlemr.Changes)
			sizes = append(sizes, size)
		}
	}

	sort.Slice(sizes, func(i, j int) bool {
		return sizes[i].Date.Before(sizes[j].Date)
	})

	return sizes, err
}

func (m *mergeRequestSizes) calculateChanges(createdat time.Time, changes []*models.MergeRequestChanges) models.ItemCount {
	newline := 0
	deletedline := 0
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

	return models.ItemCount{
		Date:  createdat,
		Count: float64(newline + deletedline),
	}
}
