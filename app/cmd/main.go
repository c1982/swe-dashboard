package main

import (
	"fmt"
	"os"
	"swe-dashboard/internal/metrics/mergerequestcomments"
	"swe-dashboard/internal/pusher/victoriametrics"
	"swe-dashboard/internal/scm/gitlab"
)

func main() {
	baseURL := os.Getenv("SWE_DASHBOARD_GITLAB_BASEURL")
	token := os.Getenv("SWE_DASHBOARD_GITLAB_TOKEN")

	gitlab, err := gitlab.NewSCM(gitlab.GitlabBaseURL(baseURL), gitlab.GitlabToken(token))
	if err != nil {
		panic(err)
	}

	service := mergerequestcomments.NewMergeRequestCommentsService(gitlab)
	leaderboard, err := service.CommentsLeaderBoard("merged", "all", 10)
	if err != nil {
		panic(err)
	}

	pusher, err := victoriametrics.NewPusher(victoriametrics.SetPushURL("http://localhost:8428/api/v1/import/prometheus"))
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(leaderboard); i++ {
		payload := fmt.Sprintf(`comments_leaderboard{id="%d",user="%s", name="%s"} %f`,
			leaderboard[i].ID,
			leaderboard[i].Username,
			leaderboard[i].Name,
			leaderboard[i].Count)
		err := pusher.Push(payload)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
