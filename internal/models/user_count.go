package models

import "time"

type UserCount struct {
	ID       int
	Name     string
	Username string
	Count    float64
	Date     time.Time
}

//TODO: may be?
func (u UserCount) ToPrometheusMetric(name string) string {
	return ""
}
