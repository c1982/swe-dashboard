package longrunningmergerequests

import (
	"sort"
	"strings"
	"swe-dashboard/internal/models"
	"time"
)

const (
	tenDays    = 864000
	ninetydays = 90
)

type SCM interface {
	ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error)
	ListMergeRequestNotes(projectID int, mergeRequestID int) (comments []*models.Comment, err error)
	GetRepository(projectID int) (repository models.Repo, err error)
}

type LongRunningMergerequestsService interface {
	List() ([]models.ItemCount, error)
}

type longRunningMergerequests struct {
	scm SCM
}

func NewLongRunningMergerequestsService(scm SCM) LongRunningMergerequestsService {
	mrc := &longRunningMergerequests{
		scm: scm,
	}
	return mrc
}

func (l *longRunningMergerequests) List() (longrunnings []models.ItemCount, err error) {
	mergerequests, err := l.scm.ListMergeRequest("opened", "all", 90)
	if err != nil {
		return longrunnings, err
	}

	longrunnings = []models.ItemCount{}
	repositories := mergerequests.GroupByRepositories()
	for i := 0; i < len(repositories); i++ {
		repo, err := l.scm.GetRepository(repositories[i].ID)
		if err != nil {
			return longrunnings, err
		}

		repo.MRs = repositories[i].MRs
		for n := 0; n < len(repo.MRs); n++ {
			mr := repo.MRs[n]
			if mr.Draft {
				continue
			}

			comments, err := l.scm.ListMergeRequestNotes(repo.ID, mr.IID)
			if err != nil {
				return longrunnings, err
			}

			lastactivitycomment := l.mergeRequestLastActivity(comments)
			if lastactivitycomment == nil {
				continue
			}

			worktime := time.Now().Unix() - lastactivitycomment.CreatedAt.Unix()
			if worktime < tenDays {
				continue
			}

			longrunnings = append(longrunnings, models.ItemCount{
				Name:  repo.Name,
				Name1: l.cleanTitle(mr.Title),
				Count: float64(worktime),
			})
		}
	}

	return longrunnings, nil
}

func (l *longRunningMergerequests) mergeRequestLastActivity(comments []*models.Comment) *models.Comment {
	if comments == nil {
		return &models.Comment{}
	}

	if len(comments) < 1 {
		return &models.Comment{}
	}

	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt.After(comments[j].CreatedAt)
	})

	return comments[0]
}

func (l *longRunningMergerequests) cleanTitle(title string) string {
	title = strings.ReplaceAll(title, "\"", "")
	title = strings.ReplaceAll(title, "/", "-")
	title = strings.ReplaceAll(title, "{", "")
	title = strings.ReplaceAll(title, "}", "")
	return title
}
