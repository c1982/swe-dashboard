package mergerequestparticipants

import (
	"sort"
	"swe-dashboard/internal/models"
)

type engageItem struct {
	Author      *models.User
	Participant *models.User
	Count       float64
}

type SCM interface {
	ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error)
	GetMergeRequestParticipants(projectID int, mergeRequestID int) (users []*models.User, err error)
}

type MergeRequestParticipantsService interface {
	List() (users []models.UserCount, err error)
}

func NewMergeRequestParticipantsService(scm SCM) MergeRequestParticipantsService {
	return &mergeRequestParticipants{
		scm: scm,
	}
}

type mergeRequestParticipants struct {
	scm SCM
}

func (mrp *mergeRequestParticipants) List() (users []models.UserCount, err error) {
	mergerequests, err := mrp.scm.ListMergeRequest("merged", "all", 30)
	if err != nil {
		return users, err
	}

	engagements := map[int]map[int]*engageItem{}

	stats := map[int]models.UserCount{}
	for i := 0; i < len(mergerequests); i++ {
		mr := mergerequests[i]
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

		v, ok := engagements[mr.Author.ID]
		if !ok {
			participantscounts := map[int]*engageItem{}
			for p := 0; p < len(participants); p++ {
				participant := participants[p]
				pv, ok := participantscounts[participant.ID]
				if !ok {
					participantscounts[participant.ID] = &engageItem{
						Author:      mr.Author,
						Participant: participant,
						Count:       1,
					}
				} else {
					pv.Count = pv.Count + 1
					participantscounts[participant.ID] = pv
				}
			}
			engagements[mr.Author.ID] = participantscounts
		} else {
			for p := 0; p < len(participants); p++ {
				participant := participants[p]
				pv, ok := v[participant.ID]
				if !ok {
					v[participant.ID] = &engageItem{
						Author:      mr.Author,
						Participant: participant,
						Count:       1,
					}
				} else {
					pv.Count = pv.Count + 1
					v[participant.ID] = pv
				}
			}
			engagements[mr.Author.ID] = v
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
