package cycletime

import (
	"sort"
	"strings"
	"swe-dashboard/internal/models"
)

type SCM interface {
	ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error)
	GetMergeRequestCommits(projectID, mergeRequestID int) (commits []*models.Commit, err error)
	GetRepository(projectID int) (repository models.Repo, err error)
	ListMergeRequestNotes(projectID int, mergeRequestID int) (comments []*models.Comment, err error)
}

type CycleTimeService interface {
	CycleTime() ([]models.ItemCount, error)
	TimeToOpen() []models.ItemCount
	TimeToReview() []models.ItemCount
	TimeToApprove() []models.ItemCount
	TimeToMerge() []models.ItemCount
}

func NewCycleTimeService(scm SCM) CycleTimeService {
	c := &cycleTime{
		scm:           scm,
		timetoreviews: []models.ItemCount{},
	}
	return c
}

type cycleTime struct {
	scm           SCM
	timetoreviews []models.ItemCount
	timetoopens   []models.ItemCount
	timetoapprove []models.ItemCount
	timetomerge   []models.ItemCount
}

func (c *cycleTime) CycleTime() (cycletimes []models.ItemCount, err error) {
	cycletimes = []models.ItemCount{}
	mergerequests, err := c.scm.ListMergeRequest("merged", "all", 30)
	if err != nil {
		return cycletimes, err
	}

	repositories := mergerequests.GroupByRepositories()
	for i := 0; i < len(repositories); i++ {
		r := repositories[i]
		repo, err := c.scm.GetRepository(r.ID)
		if err != nil {
			return cycletimes, err
		}

		for n := 0; n < len(r.MRs); n++ {
			mr := r.MRs[n]
			commits, err := c.scm.GetMergeRequestCommits(mr.ProjectID, mr.IID)
			if err != nil {
				return cycletimes, err
			}

			comments, err := c.scm.ListMergeRequestNotes(mr.ProjectID, mr.IID)
			if err != nil {
				return cycletimes, err
			}

			mergerequestfirstcommit := c.mergeRequestFirstCommit(commits)
			if mergerequestfirstcommit == nil {
				continue
			}

			mergerequestfirstcomment := c.mergeRequestFirstComment(comments)
			if mergerequestfirstcomment == nil {
				continue
			}

			mergerequestopentime := mr.CreatedAt.Unix()
			opentime := mergerequestopentime - mergerequestfirstcommit.CreatedAt.Unix()
			c.timetoopens = append(c.timetoopens, models.ItemCount{
				Name:  repo.Name,
				Name1: c.cleanTitle(mr.Title),
				Count: float64(opentime),
			})

			timetoreview := mergerequestfirstcomment.CreatedAt.Unix() - mergerequestopentime
			c.timetoreviews = append(c.timetoreviews, models.ItemCount{
				Name:  repo.Name,
				Name1: c.cleanTitle(mr.Title),
				Count: float64(timetoreview),
			})

			mergerequestapprovalcomment := c.mergeRequestApprovalComment(comments)
			timetoapprove := mergerequestfirstcomment.CreatedAt.Unix() - mergerequestapprovalcomment.CreatedAt.Unix()
			c.timetoapprove = append(c.timetoapprove, models.ItemCount{
				Name:  repo.Name,
				Name1: c.cleanTitle(mr.Title),
				Count: float64(timetoapprove),
			})

			mergetime := mr.MergedAt.Unix() - mergerequestapprovalcomment.CreatedAt.Unix()
			c.timetomerge = append(c.timetomerge, models.ItemCount{
				Name:  repo.Name,
				Name1: c.cleanTitle(mr.Title),
				Count: float64(mergetime),
			})

			cycletime := mr.MergedAt.Unix() - mergerequestfirstcommit.CreatedAt.Unix()
			cycletimes = append(cycletimes, models.ItemCount{
				Name:  repo.Name,
				Name1: c.cleanTitle(mr.Title),
				Count: float64(cycletime),
			})
		}
	}

	return cycletimes, nil
}

//TimeToOpen Time to open (from the first commit to open)
func (c *cycleTime) TimeToOpen() []models.ItemCount {
	return c.timetoopens
}

//TimeToReview Time waiting for review (from open to the first comment)
func (c *cycleTime) TimeToReview() []models.ItemCount {
	return c.timetoreviews
}

//TimeToApprove Time to approve (from the first comment to approved)
func (c *cycleTime) TimeToApprove() []models.ItemCount {
	return c.timetoapprove
}

//TimeToMerge Time to merge (from approved to merge)
func (c *cycleTime) TimeToMerge() []models.ItemCount {
	return c.timetomerge
}

func (c *cycleTime) mergeRequestFirstCommit(commits []*models.Commit) *models.Commit {
	sort.Slice(commits, func(i, j int) bool {
		return commits[i].CreatedAt.Before(commits[j].CreatedAt)
	})
	return commits[0]
}

func (c cycleTime) mergeRequestFirstComment(comments []*models.Comment) *models.Comment {
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt.Before(comments[j].CreatedAt)
	})

	//filter organic comments
	commentIndex := 0
	for i := 0; i < len(comments); i++ {
		c := comments[i]
		if c.System {
			continue
		}

		commentIndex = i
		break
	}

	if len(comments) < 1 {
		return nil
	}

	return comments[commentIndex]
}

func (c cycleTime) mergeRequestApprovalComment(comments []*models.Comment) *models.Comment {
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt.Before(comments[j].CreatedAt)
	})

	commentIndex := 0
	for i := 0; i < len(comments); i++ {
		c := comments[i]
		if !c.System {
			continue
		}

		if !c.ApprovedNote {
			continue
		}

		commentIndex = i
		break
	}

	return comments[commentIndex]
}

func (c *cycleTime) cleanTitle(title string) string {
	title = strings.ReplaceAll(title, "\"", "")
	title = strings.ReplaceAll(title, "/", "-")
	title = strings.ReplaceAll(title, "{", "")
	title = strings.ReplaceAll(title, "}", "")
	return title
}
