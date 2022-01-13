package main

import (
	"flag"
	"swe-dashboard/internal/pusher/victoriametrics"
	"swe-dashboard/internal/scm/gitlab"
	"time"
)

func main() {
	gitlabBaseURL := flag.String("scm-gitlab-baseurl", "", "--scm-gitlab-baseurl=https://your-domain-name/api/v4")
	gitlabToken := flag.String("scm-gitlab-token", "", "--scm-gitlab-token=VERY-SECRED-TOKEN (take from gitlab access token")
	victoriametricsImportURL := flag.String("victoriametrics-importurl", "", "--victoriametrics-importurl=http://localhost:8428/api/v1/import/prometheus")
	checkInterval := flag.String("check-interval", "5s", "--check-interval=6h Valid time units are ns,us,ms,s,m,h")
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

func execute(interval time.Duration, scm *gitlab.SCM, pusher *victoriametrics.Pusher) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		//TODO: execute with run group or function of map
	}
}
