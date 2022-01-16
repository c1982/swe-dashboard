package mergerequestthroughput

import (
	"swe-dashboard/internal/models"
	"time"
)

type SCM interface {
	ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error)
	GetRepository(projectID int) (repository models.Repo, err error)
}

type MergeRequestThroughputService interface {
	Throughput() (unrevieweds []models.ItemCount, err error)
}

func NewMergeRequestThroughputService(scm SCM) MergeRequestThroughputService {
	return &mergeRequestThroughput{
		scm: scm,
	}
}

type mergeRequestThroughput struct {
	scm SCM
}

func (t *mergeRequestThroughput) Throughput() (throughputs []models.ItemCount, err error) {
	mergerequests, err := t.scm.ListMergeRequest("merged", "all", time.Now().Day())
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
		count := len(repo.MRs)
		throughputs = append(throughputs, models.ItemCount{
			Name:  repo.Name,
			Count: float64(count),
		})
	}

	return throughputs, nil
}
