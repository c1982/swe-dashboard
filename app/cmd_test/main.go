package main

import (
	"fmt"
	"os"
	"swe-dashboard/internal/metrics/activecontributors"
	"swe-dashboard/internal/metrics/cycletime"
	"swe-dashboard/internal/pusher/victoriametrics"
	"swe-dashboard/internal/scm/gitlab"
)

const (
	cycleTimeMetricName     = `cycle_time{repository="%s", title="%s"} %f`
	timeToOpenMetricName    = `time_to_open{repository="%s", title="%s"} %f`
	timetoReviewMetricName  = `time_to_review{repository="%s", title="%s"} %f`
	timetoApproveMetricName = `time_to_approve{repository="%s", title="%s"} %f`
	timetoMergeMetricName   = `time_to_merge{repository="%s", title="%s"} %f`
)

func main() {
	baseURL := os.Getenv("SWE_DASHBOARD_GITLAB_BASEURL")
	token := os.Getenv("SWE_DASHBOARD_GITLAB_TOKEN")

	gitlab, err := gitlab.NewSCM(gitlab.GitlabBaseURL(baseURL), gitlab.GitlabToken(token))
	if err != nil {
		panic(err)
	}

	pusher, err := victoriametrics.NewPusher(victoriametrics.SetPushURL("http://localhost:8428"))
	if err != nil {
		panic(err)
	}

	importContributors(gitlab, pusher)
}

func importContributors(gitlab *gitlab.SCM, p *victoriametrics.Pusher) {
	service := activecontributors.NewActiveContributors(gitlab)
	metrics, err := service.List()
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(metrics); i++ {
		payload := fmt.Sprintf(`active_contributors{repository="%s", author="%s"} %f`, metrics[i].Name, metrics[i].Name1, metrics[i].Count)
		fmt.Println(payload)
		err := p.Push(payload)
		if err != nil {
			fmt.Println(err)
		}
	}

	impact := service.Impact()
	for i := 0; i < len(impact); i++ {
		payload := fmt.Sprintf(`commit_additions{repository="%s", author="%s"} %f`, impact[i].Name, impact[i].Name1, impact[i].Count)
		fmt.Println(payload)
		err := p.Push(payload)
		if err != nil {
			fmt.Println(err)
		}
	}

	for i := 0; i < len(impact); i++ {
		payload := fmt.Sprintf(`commit_deletions{repository="%s", author="%s"} %f`, impact[i].Name, impact[i].Name1, impact[i].Count1)
		fmt.Println(payload)
		err := p.Push(payload)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func importMerics(gitlab *gitlab.SCM, p *victoriametrics.Pusher) {
	service := cycletime.NewCycleTimeService(gitlab)
	metrics, err := service.CycleTime()
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(metrics); i++ {
		payload := fmt.Sprintf(cycleTimeMetricName, metrics[i].Name, metrics[i].Name1, metrics[i].Count)
		err := p.Push(payload)
		if err != nil {
			fmt.Println(err)
		}
	}

	timetoopens := service.TimeToOpen()
	for i := 0; i < len(timetoopens); i++ {
		payload := fmt.Sprintf(timeToOpenMetricName, timetoopens[i].Name, timetoopens[i].Name1, timetoopens[i].Count)
		err := p.Push(payload)
		if err != nil {
			fmt.Println(err)
		}
	}

	timetoreview := service.TimeToReview()
	for i := 0; i < len(timetoopens); i++ {
		payload := fmt.Sprintf(timetoReviewMetricName, timetoreview[i].Name, timetoreview[i].Name1, timetoreview[i].Count)
		err := p.Push(payload)
		if err != nil {
			fmt.Println(err)
		}
	}

	timetoapprove := service.TimeToApprove()
	for i := 0; i < len(timetoopens); i++ {
		payload := fmt.Sprintf(timetoApproveMetricName, timetoapprove[i].Name, timetoapprove[i].Name1, timetoapprove[i].Count)
		err := p.Push(payload)
		if err != nil {
			fmt.Println(err)
		}
	}

	timetomerge := service.TimeToMerge()
	for i := 0; i < len(timetoopens); i++ {
		payload := fmt.Sprintf(timetoMergeMetricName, timetomerge[i].Name, timetomerge[i].Name1, timetomerge[i].Count)
		err := p.Push(payload)
		if err != nil {
			fmt.Println(err)
		}
	}
}
