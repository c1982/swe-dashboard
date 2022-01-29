package mergerequestthroughput

import (
	"swe-dashboard/internal/models"
)

type SCM interface {
	ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error)
	GetRepository(projectID int) (repository models.Repo, err error)
}

type MergeRequestThroughputService interface {
	List() (unrevieweds []models.ItemCount, err error)
}

func NewMergeRequestThroughputService(scm SCM) MergeRequestThroughputService {
	return &mergeRequestThroughput{
		scm: scm,
	}
}

type mergeRequestThroughput struct {
	scm SCM
}

func (t *mergeRequestThroughput) List() (throughputs []models.ItemCount, err error) {
	mergerequests, err := t.scm.ListMergeRequest("merged", "all", 30)
	if err != nil {
		return throughputs, err
	}

	throughputs = []models.ItemCount{}
	repositories := mergerequests.GroupByRepositories()
	for i := 0; i < len(repositories); i++ {
		repo, err := t.scm.GetRepository(repositories[i].ID)
		if err != nil {
			return throughputs, err
		}

		repo.MRs = repositories[i].MRs
		days := repo.MRs.CountByDay()
		for d := 0; d < len(days); d++ {
			day := days[d]
			throughputs = append(throughputs, models.ItemCount{
				Name:  repo.Name,
				Date:  day.Date,
				Count: day.Count,
			})
		}
	}

	return throughputs, nil
}
