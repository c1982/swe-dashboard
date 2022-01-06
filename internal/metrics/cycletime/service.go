package cycletime

import (
	"sort"
	"swe-dashboard/internal/models"
	"time"
)

type SCM interface {
	ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error)
	GetMergeRequestCommits(projectID, mergeRequestID int) (commits []*models.Commit, err error)
	GetRepository(projectID int) (repository models.Repo, err error)
}

type CycleTimeService interface {
	TimeToOpen() error
	TimeToReview()
	TimeToApprove()
	TimeToMerge()
}

type cycleTime struct {
	scm SCM
}

//TimeToOpen Time to open (from the first commit to open)
func (c *cycleTime) TimeToOpen() (opentimes []models.ItemCount, err error) {
	opentimes = []models.ItemCount{}
	mergerequests, err := c.scm.ListMergeRequest("merged", "all", time.Now().Day())
	if err != nil {
		return opentimes, err
	}

	repositories := mergerequests.GroupByRepositories()
	for i := 0; i < len(repositories); i++ {
		r := repositories[i]
		repo, err := c.scm.GetRepository(r.ID)
		if err != nil {
			return opentimes, err
		}

		for n := 0; n < len(r.MRs); n++ {
			mr := r.MRs[n]
			commits, err := c.scm.GetMergeRequestCommits(mr.ProjectID, mr.IID)
			if err != nil {
				return opentimes, err
			}

			mergerequestopentime := mr.CreatedAt
			mergerequestfirstcommit := c.mergeRequestFirstCommit(commits)
			opentime := mergerequestopentime.Sub(mergerequestfirstcommit.CreatedAt)
			opentimes = append(opentimes, models.ItemCount{
				Name:  repo.Name,
				Name1: mr.Title,
				Count: float64(opentime),
			})
		}
	}

	return opentimes, nil
}

//TimeToReview Time waiting for review (from open to the first comment)
func (c *cycleTime) TimeToReview() {

}

//TimeToApprove Time to approve (from the first comment to approved)
func (c *cycleTime) TimeToApprove() {

}

//TimeToMerge Time to merge (from approved to merge)
func (c *cycleTime) TimeToMerge() {

}

func (c *cycleTime) mergeRequestFirstCommit(commits []*models.Commit) *models.Commit {
	sort.Slice(commits, func(i, j int) bool {
		return commits[i].CreatedAt.After(commits[j].CreatedAt)
	})
	return commits[0]
}
