package main

import (
	"fmt"
	"os"
	"swe-dashboard/internal/metrics/mergerequestcomments"
	"swe-dashboard/internal/metrics/mergerequestparticipants"
	"swe-dashboard/internal/metrics/mergerequestrate"
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

	pusher, err := victoriametrics.NewPusher(victoriametrics.SetPushURL("http://localhost:8428/api/v1/import/prometheus"))
	if err != nil {
		panic(err)
	}

	//importMRCommentsLeaderBoard(gitlab, pusher)
	//importMergeRequestParticipants(gitlab, pusher)
	importMergeRequestrate(gitlab, pusher)
}

func importMRCommentsLeaderBoard(gitlab *gitlab.SCM, pusher *victoriametrics.Pusher) {
	service := mergerequestcomments.NewMergeRequestCommentsService(gitlab)
	leaderboard, err := service.CommentsLeaderBoard("merged", "all", 15)
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

func importMergeRequestParticipants(gitlab *gitlab.SCM, pusher *victoriametrics.Pusher) {
	service := mergerequestparticipants.NewMergeRequestParticipantsService(gitlab)
	leaderboard, err := service.ParticipantsLeaderBoard("merged", "all", 15)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(leaderboard); i++ {
		payload := fmt.Sprintf(`participants_leaderboard{id="%d",user="%s", name="%s"} %f`,
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

func importMergeRequestrate(gitlab *gitlab.SCM, pusher *victoriametrics.Pusher) {
	service := mergerequestrate.NewMergeRequestRateService(gitlab)
	counts, err := service.MergeRequestRates("merged", "all", 7)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(counts); i++ {
		payload := fmt.Sprintf("mergerequest_rate{} %f", counts[i].Count)
		pusher.PushWithTime(payload, counts[i].Date)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
