package main

import (
	"fmt"
	"os"
	"swe-dashboard/internal/metrics/activecontributors"
	"swe-dashboard/internal/metrics/assetiterations"
	"swe-dashboard/internal/metrics/cycletime"
	"swe-dashboard/internal/metrics/mergerequestparticipants"
	"swe-dashboard/internal/pusher/victoriametrics"
	"swe-dashboard/internal/scm/gitlab"
)

const (
	cycleTimeMetricName          = `cycle_time{repository="%s", title="%s"} %f`
	timeToOpenMetricName         = `time_to_open{repository="%s", title="%s"} %f`
	timetoReviewMetricName       = `time_to_review{repository="%s", title="%s"} %f`
	timetoApproveMetricName      = `time_to_approve{repository="%s", title="%s"} %f`
	timetoMergeMetricName        = `time_to_merge{repository="%s", title="%s"} %f`
	activeContributorsMetricName = `active_contributors{repository="%s", author="%s", email="%s"} %f`
	commitAdditionsMetricName    = `commit_additions{repository="%s", author="%s", email="%s"} %f`
	commitDeletionsMetricName    = `commit_deletions{repository="%s", author="%s", email="%s"} %f`

	mergeRequestParticipantsdMetricName          = `merge_request_participants{repository="%s",user="%s", name="%s"} %f`
	mergeRequestEngagementsMetricName            = `merge_request_engagement{repository="%s", author="%s", mergedby="%s"} %f`
	mergeRequestEngagementParticipantsMetricName = `merge_request_engage_participants{repository="%s",author="%s", participant="%s"} %f`

	assetWorkingTimesWeights    = `assets_weights{repository="%s",name="%s"} %f`
	assetWorkingTimesHours      = `assets_iteration_hours{repository="%s",name="%s"} %f`
	assetWorkingTimesIterations = `assets_iterations{repository="%s",name="%s"} %f`
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

	AssetWorkingTime(gitlab, pusher)
}

func AssetWorkingTime(gitlab *gitlab.SCM, p *victoriametrics.Pusher) {
	svc := assetiterations.NewAssetIterationTimeService(gitlab, ".jpg", ".jpeg", ".png", ".psd", ".psb", ".ai")
	err := svc.CalculateChanges()
	if err != nil {
		fmt.Println(err)
	}

	weights := svc.Weights()
	for i := 0; i < len(weights); i++ {
		w := weights[i]
		payload := fmt.Sprintf(assetWorkingTimesWeights, w.Name, w.Name1, w.Count)
		err := p.Push(payload)
		if err != nil {
			fmt.Println(err)
		}
	}

	workinghours := svc.IterationHours()
	for i := 0; i < len(workinghours); i++ {
		w := workinghours[i]
		payload := fmt.Sprintf(assetWorkingTimesHours, w.Name, w.Name1, w.Count)
		err := p.Push(payload)
		if err != nil {
			fmt.Println(err)
		}
	}

	iterations := svc.Iterations()
	for i := 0; i < len(iterations); i++ {
		w := iterations[i]
		payload := fmt.Sprintf(assetWorkingTimesIterations, w.Name, w.Name1, w.Count)
		err := p.Push(payload)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func ImportMergerequestParticipants(gitlab *gitlab.SCM, p *victoriametrics.Pusher) {
	service := mergerequestparticipants.NewMergeRequestParticipantsService(gitlab)
	metrics, err := service.List()
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(metrics); i++ {
		payload := fmt.Sprintf(mergeRequestParticipantsdMetricName,
			metrics[i].Name,
			metrics[i].Name1,
			metrics[i].Name2,
			metrics[i].Count)
		fmt.Println(payload)
		err := p.Push(payload)
		if err != nil {
			fmt.Println(err)
		}
	}

	engageparticipants := service.EngageParticipants()
	for i := 0; i < len(engageparticipants); i++ {
		payload := fmt.Sprintf(mergeRequestEngagementParticipantsMetricName, engageparticipants[i].Name, engageparticipants[i].Name1, engageparticipants[i].Name2, engageparticipants[i].Count)
		fmt.Println(payload)
		err := p.Push(payload)
		if err != nil {
			fmt.Println(err)
		}
	}

	engagements := service.Engagements()
	for i := 0; i < len(engagements); i++ {
		payload := fmt.Sprintf(mergeRequestEngagementsMetricName, engagements[i].Name, engagements[i].Name1, engagements[i].Name2, engagements[i].Count)
		fmt.Println(payload)
		err := p.Push(payload)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func importContributors(gitlab *gitlab.SCM, p *victoriametrics.Pusher) {
	service := activecontributors.NewActiveContributors(gitlab)
	metrics, err := service.List()
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(metrics); i++ {
		payload := fmt.Sprintf(activeContributorsMetricName, metrics[i].Name, metrics[i].Name1, metrics[i].Name2, metrics[i].Count)
		fmt.Println(payload)
		err := p.Push(payload)
		if err != nil {
			fmt.Println(err)
		}
	}

	impact := service.Impact()
	for i := 0; i < len(impact); i++ {
		payload := fmt.Sprintf(commitAdditionsMetricName, impact[i].Name, impact[i].Name1, impact[i].Name2, impact[i].Count)
		fmt.Println(payload)
		err := p.Push(payload)
		if err != nil {
			fmt.Println(err)
		}
	}

	for i := 0; i < len(impact); i++ {
		payload := fmt.Sprintf(commitDeletionsMetricName, impact[i].Name, impact[i].Name1, impact[i].Name2, impact[i].Count1)
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
