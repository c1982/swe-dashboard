package main

import (
	"fmt"
	"os"
	"swe-dashboard/internal/metrics/mergerequestparticipants"
	"swe-dashboard/internal/scm/gitlab"
)

func main() {
	baseURL := os.Getenv("SWE_DASHBOARD_GITLAB_BASEURL")
	token := os.Getenv("SWE_DASHBOARD_GITLAB_TOKEN")

	gitlab, err := gitlab.NewSCM(gitlab.GitlabBaseURL(baseURL), gitlab.GitlabToken(token))
	if err != nil {
		panic(err)
	}

	service := mergerequestparticipants.NewMergeRequestParticipantsService(gitlab)
	leaderboard, err := service.MergeRequestParticipantsLeaderBoard("merged", "all", 10)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(leaderboard); i++ {
		fmt.Printf("%s\t%f\r\n", leaderboard[i].Username, leaderboard[i].Count)
	}
}
