package main

import (
	"flag"
	"swe-dashboard/internal/pusher/victoriametrics"
	"swe-dashboard/internal/scm/gitlab"
	"time"

	"github.com/rs/zerolog/log"
)

func main() {
	gitlabBaseURL := flag.String("scm-gitlab-baseurl", "", "--scm-gitlab-baseurl=https://your-domain-name/api/v4")
	gitlabToken := flag.String("scm-gitlab-token", "", "--scm-gitlab-token=VERY-SECRED-TOKEN (take from gitlab access token")
	//TODO: victoriametric basic auth params
	victoriametricsImportURL := flag.String("victoriametrics-importurl", "", "--victoriametrics-importurl=http://localhost:8428/api/v1/import/prometheus")
	checkInterval := flag.String("check-interval", "6h", "--check-interval=6h Valid time units are ns,us,ms,s,m,h")
	flag.Parse()

	log.Info().Str("scm-gitlab-baseurl", *gitlabBaseURL).Send()
	log.Info().Str("victoriametrics-importurl", *victoriametricsImportURL).Send()
	log.Info().Str("check-interval", *checkInterval).Send()

	interval, err := time.ParseDuration(*checkInterval)
	if err != nil {
		log.Fatal().Err(err).Str("interval", *checkInterval).Msg("interval could not be parsed")
		return
	}

	gitlab, err := gitlab.NewSCM(gitlab.GitlabBaseURL(*gitlabBaseURL), gitlab.GitlabToken(*gitlabToken))
	if err != nil {
		log.Fatal().Err(err).Str("url", *gitlabBaseURL).Send()
		return
	}

	pusher, err := victoriametrics.NewPusher(victoriametrics.SetPushURL(*victoriametricsImportURL))
	if err != nil {
		log.Fatal().Err(err).Str("url", *victoriametricsImportURL).Send()
		return
	}

	log.Info().Msg("running swed daemon")
	run(interval, gitlab, pusher)
}
