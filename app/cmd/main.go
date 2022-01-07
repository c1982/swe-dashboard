package main

import (
	"fmt"
	"os"
	"swe-dashboard/internal/metrics/cycletime"
	"swe-dashboard/internal/metrics/mergerequestcomments"
	"swe-dashboard/internal/metrics/mergerequestparticipants"
	"swe-dashboard/internal/metrics/mergerequestrate"
	"swe-dashboard/internal/metrics/mergerequestsize"
	"swe-dashboard/internal/metrics/selfmerging"
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
	//importMergeRequestrate(gitlab, pusher)
	//importMergeRequestSize(gitlab, pusher)
	//importSelfMergingUsers(gitlab, pusher)
	imporCycleTime(gitlab, pusher)
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
	counts, err := service.MergeRequestRatesThisMonth()
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(counts); i++ {
		payload := fmt.Sprintf(`merge_request_rate{repository="%s"} %f`, counts[i].Name, counts[i].Count)
		fmt.Println(payload)
		err := pusher.Push(payload)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func importMergeRequestSize(gitlab *gitlab.SCM, pusher *victoriametrics.Pusher) {
	service := mergerequestsize.NewMergeRequestSizeService(gitlab)
	counts, err := service.MergeRequestSizes()
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(counts); i++ {
		payload := fmt.Sprintf(`merge_request_size{repository="%s", title="%s"} %f`, counts[i].Name, counts[i].Name1, counts[i].Count)
		fmt.Println(payload)
		err := pusher.Push(payload)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func importSelfMergingUsers(gitlab *gitlab.SCM, pusher *victoriametrics.Pusher) {
	service := selfmerging.NewSelfMergingService(gitlab)
	users, err := service.GetSelfMergingUsers()
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(users); i++ {
		payload := fmt.Sprintf(`self_merging{name="%s", username="%s"} %f`, users[i].Name, users[i].Username, users[i].Count)
		fmt.Println(payload)
		err := pusher.Push(payload)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func imporCycleTime(gitlab *gitlab.SCM, pusher *victoriametrics.Pusher) {
	service := cycletime.NewCycleTimeService(gitlab)
	counts, err := service.CycleTime()
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(counts); i++ {
		payload := fmt.Sprintf(`cycle_time{repository="%s", title="%s"} %f`, counts[i].Name, counts[i].Name1, counts[i].Count)
		fmt.Println(payload)
		err := pusher.Push(payload)
		if err != nil {
			fmt.Println(err)
			break
		}
	}

	timetoopens := service.TimeToOpen()
	for i := 0; i < len(timetoopens); i++ {
		payload := fmt.Sprintf(`time_to_open{repository="%s", title="%s"} %f`, timetoopens[i].Name, timetoopens[i].Name1, timetoopens[i].Count)
		fmt.Println(payload)
		err := pusher.Push(payload)
		if err != nil {
			fmt.Println(err)
			break
		}
	}

	timetoreview := service.TimeToReview()
	for i := 0; i < len(timetoopens); i++ {
		payload := fmt.Sprintf(`time_to_review{repository="%s", title="%s"} %f`, timetoreview[i].Name, timetoreview[i].Name1, timetoreview[i].Count)
		fmt.Println(payload)
		err := pusher.Push(payload)
		if err != nil {
			fmt.Println(err)
			break
		}
	}

	timetoapprove := service.TimeToApprove()
	for i := 0; i < len(timetoopens); i++ {
		payload := fmt.Sprintf(`time_to_approve{repository="%s", title="%s"} %f`, timetoapprove[i].Name, timetoapprove[i].Name1, timetoapprove[i].Count)
		fmt.Println(payload)
		err := pusher.Push(payload)
		if err != nil {
			fmt.Println(err)
			break
		}
	}

	timetomerge := service.TimeToMerge()
	for i := 0; i < len(timetoopens); i++ {
		payload := fmt.Sprintf(`time_to_merge{repository="%s", title="%s"} %f`, timetomerge[i].Name, timetomerge[i].Name1, timetomerge[i].Count)
		fmt.Println(payload)
		err := pusher.Push(payload)
		if err != nil {
			fmt.Println(err)
			break
		}
	}

}
