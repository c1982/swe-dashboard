package unreviewedmergerequests

import (
	"swe-dashboard/internal/models"
	"time"
)

type SCM interface {
	ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error)
	ListMergeRequestNotes(projectID int, mergeRequestID int) (comments []*models.Comment, err error)
	GetRepository(projectID int) (repository models.Repo, err error)
}

type UnreviewedMergeRequestsService interface {
	List() (unrevieweds []models.ItemCount, err error)
}

func NewUnreviewedMergerequests(scm SCM) UnreviewedMergeRequestsService {
	return &unreviewedMergeRequests{
		scm: scm,
	}
}

type unreviewedMergeRequests struct {
	scm SCM
}

func (u *unreviewedMergeRequests) List() (unrevieweds []models.ItemCount, err error) {
	mergerequests, err := u.scm.ListMergeRequest("merged", "all", time.Now().Day())
	if err != nil {
		return unrevieweds, err
	}

	unrevieweds = []models.ItemCount{}
	repositories := mergerequests.GroupByRepositories()
	for i := 0; i < len(repositories); i++ {
		repo, err := u.scm.GetRepository(repositories[i].ID)
		if err != nil {
			return unrevieweds, err
		}

		repo.MRs = repositories[i].MRs
		mrcount := len(repo.MRs)
		unreviewdcount, err := u.unreviewedCount(repo.MRs)
		if err != nil {
			return unrevieweds, err
		}
		rate := (float64(mrcount) * float64(unreviewdcount)) / 100
		unrevieweds = append(unrevieweds, models.ItemCount{
			Name:  repo.Name,
			Count: rate,
		})
	}

	return unrevieweds, nil
}

func (u *unreviewedMergeRequests) unreviewedCount(mergerequests []models.MergeRequest) (count int, err error) {
	count = 0
	for n := 0; n < len(mergerequests); n++ {
		mr := mergerequests[n]
		comments, err := u.scm.ListMergeRequestNotes(mr.ProjectID, mr.IID)
		if err != nil {
			return count, err
		}

		isUnreviewed := u.isUnreviewed(comments)
		if !isUnreviewed {
			continue
		}
		count++
	}
	return count, nil
}

func (l *unreviewedMergeRequests) isUnreviewed(comments []*models.Comment) bool {
	for i := 0; i < len(comments); i++ {
		c := comments[i]
		if c.System {
			if c.ApprovedNote {
				return false
			}
		}
	}

	return true
}
