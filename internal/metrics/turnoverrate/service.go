package turnoverrate

import (
	"sort"
	"swe-dashboard/internal/models"
	"time"
)

type SCM interface {
	ListUsers() (users models.Users, err error)
}

type TurnOverrateService interface {
	TurnOverRate() ([]models.ItemCount, error)
}

type turnOverRate struct {
	scm   SCM
	users models.Users
}

func NewTurnOverRate(scm SCM) TurnOverrateService {
	return &turnOverRate{scm: scm}
}

func (tor *turnOverRate) TurnOverRate() (rates []models.ItemCount, err error) {
	rates = []models.ItemCount{}
	tor.users, err = tor.scm.ListUsers()
	if err != nil {
		return rates, err
	}

	_, firstcreateduserdate := tor.users.FirstUserCreatedTime()
	months := tor.generateTimeArray(firstcreateduserdate, time.Now())
	monthlymembers := tor.userCountAndMonthMap(months)

	for i := 0; i < len(monthlymembers); i++ {
		if i+1 > len(monthlymembers)-1 {
			break
		}

		B := float64(monthlymembers[i].Count)
		E := float64(monthlymembers[i+1].Count)
		L := B - E
		if L < 0 {
			L = 0
		}

		if E <= 0 {
			continue
		}

		AVG := (B + E) / 2
		rate := (L / AVG) * 100

		rates = append(rates, models.ItemCount{
			Date:  monthlymembers[i].Date,
			Count: rate,
		})
	}

	return rates, nil
}

func (tor *turnOverRate) generateTimeArray(from, to time.Time) []time.Time {
	dates := []time.Time{}
	for y := from.Year(); y <= to.Year(); y++ {
		starmonth := 1
		endmonth := 12

		if y == from.Year() {
			starmonth = int(from.Month())
		}

		if y == to.Year() {
			endmonth = int(to.Month())
		}

		for m := starmonth; m <= endmonth; m++ {
			dates = append(dates, time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.Now().Location()))
		}
	}

	return dates
}

func (tor *turnOverRate) userCountAndMonthMap(months []time.Time) []models.UserCount {
	users := tor.users.ToSlice()
	usercount := map[time.Time]int{}
	for i := 0; i < len(months); i++ {
		m := months[i]
		_, ok := usercount[m]
		if !ok {
			usercount[m] = 0
		}

		for u := 0; u < len(users); u++ {
			createdAt := users[u].CreatedAt
			lastactivity := users[u].LastActivityOn

			endmonth := m.AddDate(0, 1, 0)
			if createdAt.Before(endmonth) {
				usercount[m] = usercount[m] + 1
			}

			if lastactivity.Before(endmonth) {
				usercount[m] = usercount[m] - 1
			}
		}
	}

	membercount := []models.UserCount{}
	for k, v := range usercount {
		membercount = append(membercount, models.UserCount{
			Date:  k,
			Count: float64(v),
		})
	}

	sort.Slice(membercount, func(i, j int) bool {
		return membercount[i].Date.Before(membercount[j].Date)
	})

	return membercount
}
