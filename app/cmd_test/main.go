package main

import (
	"fmt"
	"os"
	"swe-dashboard/internal/metrics/mergerequestsuccessrate"
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

	importMerics(gitlab, pusher)
}

func importMerics(gitlab *gitlab.SCM, pusher *victoriametrics.Pusher) {
	service := mergerequestsuccessrate.NewMergeRequestSuccessRateService(gitlab)
	counts, err := service.List()
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < len(counts); i++ {
		payload := fmt.Sprintf(`merge_request_success_rate{repository="%s"} %f`, counts[i].Name, counts[i].Count)
		fmt.Println(payload)
		err := pusher.Push(payload)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
