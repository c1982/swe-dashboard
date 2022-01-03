package models

import (
	"testing"
	"time"
)

var (
	testUserList = Users{
		1:  User{ID: 1, CreatedAt: time.Date(2019, 2, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Date(2020, 10, 1, 0, 0, 0, 0, time.Now().Location())},
		2:  User{ID: 2, CreatedAt: time.Date(2019, 2, 1, 5, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		3:  User{ID: 3, CreatedAt: time.Date(2019, 3, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		4:  User{ID: 4, CreatedAt: time.Date(2019, 5, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Date(2019, 12, 1, 0, 0, 0, 0, time.Now().Location())},
		5:  User{ID: 5, CreatedAt: time.Date(2019, 12, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		6:  User{ID: 6, CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		7:  User{ID: 7, CreatedAt: time.Date(2020, 3, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		8:  User{ID: 8, CreatedAt: time.Date(2020, 3, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Date(2021, 2, 1, 0, 0, 0, 0, time.Now().Location())},
		9:  User{ID: 9, CreatedAt: time.Date(2020, 3, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Date(2020, 3, 1, 0, 0, 0, 0, time.Now().Location())},
		10: User{ID: 10, CreatedAt: time.Date(2020, 4, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		11: User{ID: 11, CreatedAt: time.Date(2020, 5, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		12: User{ID: 12, CreatedAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		13: User{ID: 13, CreatedAt: time.Date(2020, 11, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		14: User{ID: 14, CreatedAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Date(2021, 1, 1, 0, 0, 0, 0, time.Now().Location())},
		15: User{ID: 15, CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		16: User{ID: 16, CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		17: User{ID: 17, CreatedAt: time.Date(2021, 3, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		18: User{ID: 18, CreatedAt: time.Date(2021, 4, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		19: User{ID: 19, CreatedAt: time.Date(2021, 5, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Date(2021, 6, 1, 0, 0, 0, 0, time.Now().Location())},
		20: User{ID: 20, CreatedAt: time.Date(2021, 6, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Date(2019, 7, 1, 0, 0, 0, 0, time.Now().Location())},
		21: User{ID: 21, CreatedAt: time.Date(2021, 7, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Date(2021, 8, 1, 0, 0, 0, 0, time.Now().Location())},
		22: User{ID: 22, CreatedAt: time.Date(2021, 8, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		23: User{ID: 23, CreatedAt: time.Date(2021, 8, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		24: User{ID: 24, CreatedAt: time.Date(2021, 9, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		25: User{ID: 25, CreatedAt: time.Date(2021, 10, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		26: User{ID: 26, CreatedAt: time.Date(2021, 11, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
		27: User{ID: 27, CreatedAt: time.Date(2021, 11, 1, 0, 0, 0, 0, time.Now().Location()), LastActivityOn: time.Now()},
	}
)

func TestFirstUserCreatedTime(t *testing.T) {
	wantID := 1
	wantDate := time.Date(2019, 2, 1, 0, 0, 0, 0, time.Now().Location())
	users := testUserList
	id, createdat := users.FirstUserCreatedTime()
	if id != wantID {
		t.Errorf("invalid user ID, got: %d, want: %d", id, wantID)
	}

	if createdat.Year() != wantDate.Year() &&
		createdat.Month() != wantDate.Month() &&
		createdat.Day() != wantDate.Day() {
		t.Errorf("got: %s, want: %s", createdat, wantDate)
	}
}
