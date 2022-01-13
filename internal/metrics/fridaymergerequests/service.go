package fridaymergerequests

import (
	"swe-dashboard/internal/models"
	"time"
)

type SCM interface {
	ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error)
	GetRepository(projectID int) (repository models.Repo, err error)
}

type FridayMergerequestsService interface {
	List() (fridaymergerequests []models.ItemCount, err error)
}

func NewFridayMergeRequests(scm SCM) FridayMergerequestsService {
	return &fridayMergerequests{
		scm: scm,
	}
}

type fridayMergerequests struct {
	scm SCM
}

func (f *fridayMergerequests) List() (fridaymergerequests []models.ItemCount, err error) {
	mergerequests, err := f.scm.ListMergeRequest("", "all", time.Now().Day())
	if err != nil {
		return fridaymergerequests, err
	}

	fridaymergerequests = []models.ItemCount{}
	repositories := mergerequests.GroupByRepositories()
	for i := 0; i < len(repositories); i++ {
		repo, err := f.scm.GetRepository(repositories[i].ID)
		if err != nil {
			return fridaymergerequests, err
		}
		repo.MRs = repositories[i].MRs
		count := f.fridayOpenedMergeRequests(repo.MRs)
		fridaymergerequests = append(fridaymergerequests, models.ItemCount{
			Name:  repo.Name,
			Count: float64(count),
		})
	}

	return fridaymergerequests, nil
}

func (f *fridayMergerequests) fridayOpenedMergeRequests(mergerequests []models.MergeRequest) (count int) {
	count = 0
	for n := 0; n < len(mergerequests); n++ {
		mr := mergerequests[n]
		if mr.CreatedAt.Weekday() != time.Friday {
			continue
		}

		count++
	}
	return count
}
