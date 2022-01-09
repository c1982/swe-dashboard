package longrunningmergerequests

import (
	"sort"
	"swe-dashboard/internal/models"
	"time"
)

type SCM interface {
	ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error)
	ListMergeRequestNotes(projectID int, mergeRequestID int) (comments []*models.Comment, err error)
}

type LongRunningMergerequestsService interface {
	List() ([]models.ItemCount, error)
}

type longRunningMergerequests struct {
	scm SCM
}

func NewMergeRequestCommentsService(scm SCM) LongRunningMergerequestsService {
	mrc := &longRunningMergerequests{
		scm: scm,
	}
	return mrc
}

func (l *longRunningMergerequests) List() (longrunnings []models.ItemCount, err error) {
	mergerequests, err := l.scm.ListMergeRequest("opened", "all", time.Now().Day())
	if err != nil {
		return longrunnings, err
	}

	longrunnings = []models.ItemCount{}

	//TODO: repos first!
	for i := 0; i < len(mergerequests); i++ {
		mr := mergerequests[i]
		if mr.Draft {
			continue
		}

		comments, err := l.scm.ListMergeRequestNotes(mr.ProjectID, mr.IID)
		if err != nil {
			return longrunnings, err
		}

		mergerequestcreatetime := mr.CreatedAt.Unix()
		lastactivitycomment := l.mergeRequestLastActivity(comments)
		worktime := mergerequestcreatetime - lastactivitycomment.CreatedAt.Unix()

		//TODO: 1 month is fit to long runnings
		if worktime < 100 {
			continue
		}

		longrunnings = append(longrunnings, models.ItemCount{
			Name:  "",
			Name1: mr.Title,
			Count: float64(worktime),
		})
	}

	return longrunnings, nil
}

func (l *longRunningMergerequests) mergeRequestLastActivity(comments []*models.Comment) *models.Comment {
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt.After(comments[j].CreatedAt)
	})

	return comments[0]
}
