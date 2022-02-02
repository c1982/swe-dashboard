package selfmerging

import (
	"sort"
	"swe-dashboard/internal/models"
	"time"
)

const (
	mrStateMerged = "merged"
)

type SCM interface {
	ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error)
}

type SelfMergingService interface {
	GetSelfMergingUsers() (users []models.UserCount, err error)
}

type selfMerging struct {
	scm SCM
}

func NewSelfMergingService(scm SCM) SelfMergingService {
	return &selfMerging{scm: scm}
}

func (s *selfMerging) GetSelfMergingUsers() (users []models.UserCount, err error) {
	tmpusers := map[int]models.UserCount{}
	users = []models.UserCount{}

	mrs, err := s.scm.ListMergeRequest("merged", "all", 30)
	if err != nil {
		return users, err
	}

	for i := 0; i < len(mrs); i++ {
		mr := mrs[i]
		if mr.Author.ID != mr.MergedBy.ID {
			continue
		}

		v, ok := tmpusers[mr.Author.ID]
		if !ok {
			tmpusers[mr.Author.ID] = models.UserCount{
				ID:       mr.Author.ID,
				Name:     mr.Author.Name,
				Username: mr.Author.Username,
				Count:    1,
			}
		} else {
			v.Count = v.Count + 1
			tmpusers[mr.Author.ID] = v
		}
	}

	for _, v := range tmpusers {
		users = append(users, v)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].Count > users[j].Count
	})

	return users, nil
}
