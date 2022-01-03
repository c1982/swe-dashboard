package mergerequestrate

import (
	"sort"
	"swe-dashboard/internal/models"
	"time"
)

type SCM interface {
	ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error)
	ListUsers() (users models.Users, err error)
}

//MergeRequestRateService Return Pull/Merge Request rates service
type MergeRequestRateService interface {
	/*MergeRequestRates calculate merge request rates
	state: opened, closed, locked, merged.
	scope: created_by_me, assigned_to_me, all
	createdafterday: Return merge requests created on or after the given day.*/
	MergeRequestRates(state, scope string, createdafterday int) ([]models.ItemCount, error)
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

func (m *mergeRequestRates) MergeRequestRates(state, scope string, createdafterday int) (rates []models.ItemCount, err error) {
	m.members, err = m.scm.ListUsers()
	if err != nil {
		return rates, err
	}

	m.mergerequests, err = m.scm.ListMergeRequest(state, scope, createdafterday)
	if err != nil {
		return rates, err
	}

	team := m.groupMemberByTime()
	mrcountspermonth := m.mergerequests.CountByMonth()
	rates = []models.ItemCount{}

	for i := 0; i < len(mrcountspermonth); i++ {
		mr := mrcountspermonth[i]
		members, ok := team[mr.Date]
		if !ok {
			members = models.ItemCount{
				Date:  mr.Date,
				Count: float64(len(team)),
			}
		}

		mrrate := 0.0
		if mr.Count > 0 && members.Count > 0 {
			mrrate = float64(mr.Count) / float64(members.Count)
		}
		rates = append(rates, models.ItemCount{
			Date:  mr.Date,
			Count: mrrate})
	}

	return
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
