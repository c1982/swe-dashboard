package models

import (
	"sort"
	"time"
)

const (
	mrStateMerged = "merged"
)

type MergeRequests []MergeRequest

func (m MergeRequests) CountByMonth() []ItemCount {
	list := []ItemCount{}
	groupmonth := map[time.Time]int{}
	for i := 0; i < len(m); i++ {
		mr := m[i]
		month := time.Date(mr.CreatedAt.Year(), mr.CreatedAt.Month(), 1, 0, 0, 0, 0, time.Now().Location())
		v, ok := groupmonth[month]
		if ok {
			groupmonth[month] = v + 1
		} else {
			groupmonth[month] = 1
		}
	}

	for date, count := range groupmonth {
		list = append(list, ItemCount{
			Date:  date,
			Count: float64(count),
		})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Date.Before(list[j].Date)
	})

	return list

}

func (m MergeRequests) GroupByRepositories() []*Repo {
	group := map[int]*Repo{}
	for i := 0; i < len(m); i++ {
		mr := m[i]
		repo, ok := group[mr.ProjectID]
		if !ok {
			mrs := []MergeRequest{}
			mrs = append(mrs, mr)
			group[mr.ProjectID] = &Repo{
				ID:  mr.ProjectID,
				MRs: mrs,
			}
		} else {
			repo.MRs = append(repo.MRs, mr)
		}
	}

	repos := []*Repo{}
	for _, r := range group {
		repos = append(repos, r)
	}

	return repos
}
