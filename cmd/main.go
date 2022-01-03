package main

import (
	"os"
	"swe-dashboard/internal/metrics/selfmerging"
	"swe-dashboard/internal/metrics/turnoverrate"
	"swe-dashboard/internal/scm/gitlab"
)

func main() {
	baseURL := os.Getenv("SWE_DASHBOARD_GITLAB_BASEURL")
	token := os.Getenv("SWE_DASHBOARD_GITLAB_TOKEN")
	gitlab, err := gitlab.NewSCM(gitlab.GitlabBaseURL(baseURL), gitlab.GitlabToken(token))
	if err != nil {
		panic(err)
	}

	rate := turnoverrate.NewTurnOverRate(gitlab)
	rate.Calculate()

	sm := selfmerging.NewSelfMergingService(gitlab)
	sm.Calculate()

}
