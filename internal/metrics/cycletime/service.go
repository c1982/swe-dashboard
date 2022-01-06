package cycletime

import (
	"sort"
	"strings"
	"swe-dashboard/internal/models"
	"time"
)

type SCM interface {
	ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error)
	GetMergeRequestCommits(projectID, mergeRequestID int) (commits []*models.Commit, err error)
	GetRepository(projectID int) (repository models.Repo, err error)
}

type CycleTimeService interface {
	TimeToOpen() ([]models.ItemCount, error)
	TimeToReview()
	TimeToApprove()
	TimeToMerge()
}

func NewCycleTimeService(scm SCM) CycleTimeService {
	c := &cycleTime{
		scm: scm,
	}

	return c
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

			mergerequestopentime := mr.CreatedAt.Unix()
			mergerequestfirstcommit := c.mergeRequestFirstCommit(commits)
			opentime := mergerequestopentime - mergerequestfirstcommit.CreatedAt.Unix()
			opentimes = append(opentimes, models.ItemCount{
				Name:  repo.Name,
				Name1: c.cleanTitle(mr.Title),
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
		return commits[i].CreatedAt.Before(commits[j].CreatedAt)
	})
	return commits[0]
}

func (c *cycleTime) cleanTitle(title string) string {
	return strings.TrimFunc(title, func(r rune) bool {
		switch r {
		case '}', '"', '{':
			return true
		}

		return false
	})
}
