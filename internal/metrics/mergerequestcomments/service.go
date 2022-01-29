package mergerequestcomments

import (
	"sort"
	"swe-dashboard/internal/models"
)

const (
	mrStateMerged = "merged"
)

type SCM interface {
	ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error)
	ListMergeRequestNotes(projectID int, mergeRequestID int) (comments []*models.Comment, err error)
}

type MergeRequestCommentsService interface {
	List() (users []models.UserCount, err error)
}

type mergeRequestComments struct {
	scm SCM
}

func NewMergeRequestCommentsService(scm SCM) MergeRequestCommentsService {
	mrc := &mergeRequestComments{
		scm: scm,
	}
	return mrc
}

func (mrc *mergeRequestComments) List() (users []models.UserCount, err error) {
	stats := map[int]models.UserCount{}
	mergerequests, err := mrc.scm.ListMergeRequest("merged", "all", 30)
	if err != nil {
		return users, err
	}

	for i := 0; i < len(mergerequests); i++ {
		mr := mergerequests[i]

		comments, err := mrc.scm.ListMergeRequestNotes(mr.ProjectID, mr.IID)
		if err != nil {
			return users, err
		}

		for u := 0; u < len(comments); u++ {
			comment := comments[u]
			if comment.System {
				continue
			}

			v, ok := stats[comment.Author.ID]
			if !ok {
				stats[comment.Author.ID] = models.UserCount{
					ID:       comment.Author.ID,
					Name:     comment.Author.Name,
					Username: comment.Author.Username,
					Count:    1,
				}
			} else {
				v.Count = v.Count + 1
				stats[comment.Author.ID] = v
			}
		}
	}

	users = []models.UserCount{}
	for _, v := range stats {
		users = append(users, v)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].Count > users[j].Count
	})
	return users, nil
}
