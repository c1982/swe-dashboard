package mergerequestengagement

import (
	"fmt"
	"swe-dashboard/internal/models"
)

type engageItem struct {
	Author *models.User
	Merger *models.User
	Count  float64
}

type SCM interface {
	GetRepository(projectID int) (repository models.Repo, err error)
	ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error)
}

type MergeRequestEngagementService interface {
	List() (rates []models.ItemCount, err error)
}

func NewMergeRequestEngagementService(scm SCM) MergeRequestEngagementService {
	m := &mergeEngagements{
		scm: scm,
	}
	return m
}

type mergeEngagements struct {
	scm SCM
}

func (m *mergeEngagements) List() (engagements []models.ItemCount, err error) {
	engagements = []models.ItemCount{}
	mergerequests, err := m.scm.ListMergeRequest("merged", "all", 30)
	if err != nil {
		return engagements, err
	}

	repositories := mergerequests.GroupByRepositories()
	for i := 0; i < len(repositories); i++ {
		repo, err := m.scm.GetRepository(repositories[i].ID)
		if err != nil {
			return engagements, err
		}
		repo.MRs = repositories[i].MRs
		counts := m.mergerCount(repo.MRs)

		for _, v := range counts {
			engagements = append(engagements, models.ItemCount{
				Name:  repo.Name,
				Name1: v.Author.Name,
				Name2: v.Merger.Name,
				Count: v.Count,
			})
		}
	}
	return engagements, nil
}

func (m *mergeEngagements) mergerCount(mergerequests models.MergeRequests) map[string]*engageItem {
	counts := map[string]*engageItem{}
	for n := 0; n < len(mergerequests); n++ {
		mr := mergerequests[n]
		//passed self-mergeds
		if mr.Author.ID == mr.MergedBy.ID {
			continue
		}
		pair := fmt.Sprintf("%d-%d", mr.Author.ID, mr.MergedBy.ID)
		v, ok := counts[pair]
		if !ok {
			counts[pair] = &engageItem{
				Author: mr.Author,
				Merger: mr.MergedBy,
				Count:  1,
			}
		} else {
			v.Count = v.Count + 1
			counts[pair] = v
		}
	}
	return counts
}
