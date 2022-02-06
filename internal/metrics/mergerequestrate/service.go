package mergerequestrate

import (
	"sort"
	"swe-dashboard/internal/models"
	"time"
)

type SCM interface {
	GetRepository(projectID int) (repository models.Repo, err error)
	ListAllProjectMembers(projectID int) (members []models.User, err error)
	ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error)
}

//MergeRequestRateService Return Pull/Merge Request rates service
type MergeRequestRateService interface {
	MergeRequestRates() (rates []models.ItemCount, err error)
}

func NewMergeRequestRateService(scm SCM) MergeRequestRateService {
	m := &mergeRequestRates{
		scm: scm,
	}
	return m
}

type mergeRequestRates struct {
	scm           SCM
	members       models.Users
	mergerequests models.MergeRequests
}

func (m *mergeRequestRates) MergeRequestRates() (rates []models.ItemCount, err error) {
	m.mergerequests, err = m.scm.ListMergeRequest("merged", "all", 7)
	if err != nil {
		return rates, err
	}

	rates = []models.ItemCount{}
	repositories := m.mergerequests.GroupByRepositories()
	for i := 0; i < len(repositories); i++ {
		repo, err := m.scm.GetRepository(repositories[i].ID)
		if err != nil {
			return rates, err
		}

		repo.MRs = repositories[i].MRs
		members, err := m.scm.ListAllProjectMembers(repo.ID)
		if err != nil {
			return rates, err
		}

		rate := float64(len(repo.MRs)) / float64(len(members))
		rates = append(rates, models.ItemCount{
			Name:   repo.Name,
			Count:  rate,
			Count1: float64(len(members)),
		})
	}

	return rates, nil
}

func (m *mergeRequestRates) groupMembersByActivitionTime() (membercount []models.ItemCount) {
	groupmember := m.members.CountByCreatedMonth()
	totalmember := func(groups map[time.Time]int, after time.Time) int {
		t := 0
		for k, v := range groups {
			if after.After(k) {
				t = t + v
			}
		}
		return t
	}

	membercount = []models.ItemCount{}
	for date, count := range groupmember {
		total := totalmember(groupmember, date)
		membercount = append(membercount, models.ItemCount{
			Count: float64(total + count),
			Date:  date,
		})
	}

	sort.Slice(membercount, func(i, j int) bool {
		return membercount[i].Date.Before(membercount[j].Date)
	})

	return membercount
}

func (m *mergeRequestRates) groupMemberByTime() (membertimegroups map[time.Time]models.ItemCount) {
	membercounts := m.groupMembersByActivitionTime()
	membertimegroups = map[time.Time]models.ItemCount{}
	for i := 0; i < len(membercounts); i++ {
		mc := membercounts[i]
		membertimegroups[mc.Date] = mc
	}

	return membertimegroups
}
