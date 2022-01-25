package mergerequestparticipants

import (
	"sort"
	"swe-dashboard/internal/models"
)

const (
	mrStateMerged = "merged"
)

type SCM interface {
	ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error)
	GetMergeRequestParticipants(projectID int, mergeRequestID int) (users []*models.User, err error)
}

type MergeRequestParticipantsService interface {
	ParticipantsLeaderBoard() (users []models.UserCount, err error)
}

func NewMergeRequestParticipantsService(scm SCM) MergeRequestParticipantsService {
	return &mergeRequestParticipants{
		scm: scm,
	}
}

type mergeRequestParticipants struct {
	scm SCM
}

func (mrp *mergeRequestParticipants) ParticipantsLeaderBoard() (users []models.UserCount, err error) {
	mergerequests, err := mrp.scm.ListMergeRequest("merged", "all", 1)
	if err != nil {
		return users, err
	}

	stats := map[int]models.UserCount{}
	for i := 0; i < len(mergerequests); i++ {
		mr := mergerequests[i]
		if mr.State != mrStateMerged {
			continue
		}

		participants, err := mrp.scm.GetMergeRequestParticipants(mr.ProjectID, mr.IID)
		if err != nil {
			return users, err
		}

		for u := 0; u < len(participants); u++ {
			user := participants[u]
			v, ok := stats[user.ID]
			if !ok {
				stats[user.ID] = models.UserCount{
					ID:       user.ID,
					Name:     user.Name,
					Username: user.Username,
					Count:    1,
				}
			} else {
				v.Count = v.Count + 1
				stats[user.ID] = v
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
