package turnoverrate

import (
	"swe-dashboard/internal/models"
	"testing"
	"time"
)

var (
	testUserList = models.Users{
		1:  models.User{ID: 1, CreatedAt: time.Date(2019, 2, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Date(2020, 10, 1, 0, 0, 0, 0, time.Now().Location())},
		2:  models.User{ID: 2, CreatedAt: time.Date(2019, 2, 1, 5, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		3:  models.User{ID: 3, CreatedAt: time.Date(2019, 3, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		4:  models.User{ID: 4, CreatedAt: time.Date(2019, 5, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Date(2019, 12, 1, 0, 0, 0, 0, time.Now().Location())},
		5:  models.User{ID: 5, CreatedAt: time.Date(2019, 12, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		6:  models.User{ID: 6, CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		7:  models.User{ID: 7, CreatedAt: time.Date(2020, 3, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		8:  models.User{ID: 8, CreatedAt: time.Date(2020, 3, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Date(2021, 2, 1, 0, 0, 0, 0, time.Now().Location())},
		9:  models.User{ID: 9, CreatedAt: time.Date(2020, 3, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Date(2020, 3, 1, 0, 0, 0, 0, time.Now().Location())},
		10: models.User{ID: 10, CreatedAt: time.Date(2020, 4, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		11: models.User{ID: 11, CreatedAt: time.Date(2020, 5, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		12: models.User{ID: 12, CreatedAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		13: models.User{ID: 13, CreatedAt: time.Date(2020, 11, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		14: models.User{ID: 14, CreatedAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Date(2021, 1, 1, 0, 0, 0, 0, time.Now().Location())},
		15: models.User{ID: 15, CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		16: models.User{ID: 16, CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		17: models.User{ID: 17, CreatedAt: time.Date(2021, 3, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		18: models.User{ID: 18, CreatedAt: time.Date(2021, 4, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		19: models.User{ID: 19, CreatedAt: time.Date(2021, 5, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Date(2021, 6, 1, 0, 0, 0, 0, time.Now().Location())},
		20: models.User{ID: 20, CreatedAt: time.Date(2021, 6, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Date(2019, 7, 1, 0, 0, 0, 0, time.Now().Location())},
		21: models.User{ID: 21, CreatedAt: time.Date(2021, 7, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Date(2021, 8, 1, 0, 0, 0, 0, time.Now().Location())},
		22: models.User{ID: 22, CreatedAt: time.Date(2021, 8, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		23: models.User{ID: 23, CreatedAt: time.Date(2021, 8, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		24: models.User{ID: 24, CreatedAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		25: models.User{ID: 25, CreatedAt: time.Date(2021, 10, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		26: models.User{ID: 26, CreatedAt: time.Date(2021, 11, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		27: models.User{ID: 27, CreatedAt: time.Date(2021, 11, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
	}
)

func TestGenerateTimeArray(t *testing.T) {
	monthCount := 34
	provider := turnOverRate{}
	firstCreateAt := time.Date(2019, 02, 1, 0, 0, 0, 0, time.Now().Location())
	today := time.Date(2021, 11, 1, 0, 0, 0, 0, time.Now().Location())
	timelist := provider.generateTimeArray(firstCreateAt, today)

	if len(timelist) != monthCount {
		t.Errorf("unexpected months count got: %d, want: %d", len(timelist), monthCount)
	}

	if timelist[0] != firstCreateAt {
		t.Errorf("invalid first date. got: %s, want: %s", timelist[0], firstCreateAt)
	}

	if timelist[len(timelist)-1] != today {
		t.Errorf("invalid last date. got: %s, want: %s", timelist[len(timelist)-1], today)
	}
}

func TestExplodeUsersCreateTimeToMonths(t *testing.T) {
	totalmonths := 34
	expireduser := 8
	provider := turnOverRate{
		users: testUserList,
	}

	firstCreatedAt := time.Date(2019, 02, 1, 0, 0, 0, 0, time.Now().Location())
	today := time.Date(2021, 11, 1, 0, 0, 0, 0, time.Now().Location())
	months := provider.generateTimeArray(firstCreatedAt, today)

	monthsandmembercount := provider.userCountAndMonthMap(months)
	if len(monthsandmembercount) != totalmonths {
		t.Errorf("unexpected user count. got: %d, want :%d", len(monthsandmembercount), totalmonths)
	}

	firstdate := monthsandmembercount[0].Date
	lastdate := monthsandmembercount[len(monthsandmembercount)-1].Date
	lastusercount := monthsandmembercount[len(monthsandmembercount)-1].Count

	if !firstdate.Equal(firstCreatedAt) {
		t.Errorf("invalid first created at. got: %s, want :%s", firstdate, firstCreatedAt)
	}

	if !lastdate.Equal(today) {
		t.Errorf("invalid today at. got: %s, want :%s", lastdate, today)
	}

	if (len(provider.users) - expireduser) != int(lastusercount) {
		t.Errorf("expired user wrong calculated. got: %f, want: %d", lastusercount, len(provider.users)-expireduser)
	}
}
