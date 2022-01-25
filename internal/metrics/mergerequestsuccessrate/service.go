package mergerequestsuccessrate

import (
	"swe-dashboard/internal/models"
)

const (
	mrStateMerged = "merged"
	mrStateOpened = "opened"
	mrStateClosed = "closed"
	mrStateLocked = "locked"
)

type SCM interface {
	GetRepository(projectID int) (repository models.Repo, err error)
	ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error)
}

type MergeRequestSuccessRateService interface {
	List() (sizes []models.ItemCount, err error)
}

type successRate struct {
	scm SCM
}

func NewMergeRequestSuccessRateService(scm SCM) MergeRequestSuccessRateService {
	m := &successRate{
		scm: scm,
	}
	return m
}

func (s *successRate) List() (rates []models.ItemCount, err error) {
	mergerequests, err := s.scm.ListMergeRequest("", "all", 1)
	if err != nil {
		return rates, err
	}

	rates = []models.ItemCount{}
	repositories := mergerequests.GroupByRepositories()
	for i := 0; i < len(repositories); i++ {
		repo, err := s.scm.GetRepository(repositories[i].ID)
		if err != nil {
			return rates, err
		}
		repo.MRs = repositories[i].MRs
		mergedcount := 0
		closedcount := 0

		for n := 0; n < len(repo.MRs); n++ {
			mr := repo.MRs[n]
			switch mr.State {
			case mrStateClosed, mrStateLocked:
				closedcount++
			case mrStateMerged:
				mergedcount++
			}
		}

		rate := float64(mergedcount) * 100 / float64((closedcount + mergedcount))
		rates = append(rates, models.ItemCount{
			Name:  repo.Name,
			Count: float64(rate),
		})
	}

	return rates, err
}
