package main

import (
	"flag"
	"swe-dashboard/internal/metrics/cycletime"
	"swe-dashboard/internal/metrics/fridaymergerequests"
	"swe-dashboard/internal/metrics/longrunningmergerequests"
	"swe-dashboard/internal/metrics/mergerequestcomments"
	"swe-dashboard/internal/metrics/mergerequestparticipants"
	"swe-dashboard/internal/metrics/mergerequestrate"
	"swe-dashboard/internal/metrics/mergerequestsize"
	"swe-dashboard/internal/metrics/selfmerging"
	"swe-dashboard/internal/metrics/unreviewedmergerequests"
	"swe-dashboard/internal/pusher/victoriametrics"
	"swe-dashboard/internal/scm/gitlab"
	"time"
)

func main() {
	gitlabBaseURL := flag.String("scm-gitlab-baseurl", "", "--scm-gitlab-baseurl=https://your-domain-name/api/v4")
	gitlabToken := flag.String("scm-gitlab-token", "", "--scm-gitlab-token=VERY-SECRED-TOKEN (take from gitlab access token")
	victoriametricsImportURL := flag.String("victoriametrics-importurl", "", "--victoriametrics-importurl=http://localhost:8428/api/v1/import/prometheus")
	checkInterval := flag.String("check-interval", "6h", "--check-interval=6h Valid time units are ns,us,ms,s,m,h")
	flag.Parse()

	interval, err := time.ParseDuration(*checkInterval)
	if err != nil {
		flag.Usage()
		panic(err)
	}

	gitlab, err := gitlab.NewSCM(gitlab.GitlabBaseURL(*gitlabBaseURL), gitlab.GitlabToken(*gitlabToken))
	if err != nil {
		flag.Usage()
		panic(err)
	}

	pusher, err := victoriametrics.NewPusher(victoriametrics.SetPushURL(*victoriametricsImportURL))
	if err != nil {
		flag.Usage()
		panic(err)
	}

	execute(interval, gitlab, pusher)
}

func execute(interval time.Duration, gitlab *gitlab.SCM, pusher *victoriametrics.Pusher) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		_ = pusher.ImportCycleTimeMetric(cycletime.NewCycleTimeService(gitlab))
		_ = pusher.ImporFridayMergeRequests(fridaymergerequests.NewFridayMergeRequests(gitlab))
		_ = pusher.ImportLongTunningMergeRequests(longrunningmergerequests.NewLongRunningMergerequestsService(gitlab))
		_ = pusher.ImportMergeRequestComments(mergerequestcomments.NewMergeRequestCommentsService(gitlab))
		_ = pusher.ImportMergeRequestParticipants(mergerequestparticipants.NewMergeRequestParticipantsService(gitlab))
		_ = pusher.ImportMergeRequestRate(mergerequestrate.NewMergeRequestRateService(gitlab))
		_ = pusher.ImportMergeRequestSize(mergerequestsize.NewMergeRequestSizeService(gitlab))
		_ = pusher.ImportSelfMerging(selfmerging.NewSelfMergingService(gitlab))
		_ = pusher.ImportUnreviewedMergeRequests(unreviewedmergerequests.NewUnreviewedMergerequests(gitlab))
	}
}
