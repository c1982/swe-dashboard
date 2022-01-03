package models

import (
	"sort"
	"time"
)

type Users map[int]User

func (u Users) FirstUserCreatedTime() (id int, createdAt time.Time) {
	users := u.ToSlice()
	sort.Slice(users, func(i int, j int) bool {
		return users[i].CreatedAt.Before(users[j].CreatedAt)
	})
	firstUser := users[0]
	return firstUser.ID, firstUser.CreatedAt
}

func (u Users) ToSlice() []User {
	users := []User{}
	for _, v := range u {
		users = append(users, v)
	}

	return users
}

func (u Users) CountByCreatedMonth() map[time.Time]int {
	groupmember := map[time.Time]int{}
	for _, user := range u {
		month := time.Date(user.CreatedAt.Year(), user.CreatedAt.Month(), 1, 0, 0, 0, 0, time.Now().Location())
		v, ok := groupmember[month]
		if ok {
			groupmember[month] = v + 1
		} else {
			groupmember[month] = 1
		}
	}

	for _, user := range u {
		month := time.Date(user.LastActivityOn.Year(), user.LastActivityOn.Month(), 1, 0, 0, 0, 0, time.Now().Location())
		v, ok := groupmember[month]
		if ok {
			groupmember[month] = v - 1
		}
	}

	return groupmember
}
