package main

import (
	"swe-dashboard/internal/metrics/activecontributors"
	"swe-dashboard/internal/metrics/cycletime"
	"swe-dashboard/internal/metrics/defectrate"
	"swe-dashboard/internal/metrics/fridaymergerequests"
	"swe-dashboard/internal/metrics/longrunningmergerequests"
	"swe-dashboard/internal/metrics/mergerequestcomments"
	"swe-dashboard/internal/metrics/mergerequestengagement"
	"swe-dashboard/internal/metrics/mergerequestparticipants"
	"swe-dashboard/internal/metrics/mergerequestrate"
	"swe-dashboard/internal/metrics/mergerequestsize"
	"swe-dashboard/internal/metrics/mergerequestsuccessrate"
	"swe-dashboard/internal/metrics/mergerequestthroughput"
	"swe-dashboard/internal/metrics/selfmerging"
	"swe-dashboard/internal/metrics/unreviewedmergerequests"
	"swe-dashboard/internal/pusher/victoriametrics"
	"swe-dashboard/internal/scm/gitlab"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

var (
	mux = &sync.RWMutex{}
)

func setMetricsFunctions(mux *sync.RWMutex, gitlab *gitlab.SCM, pusher *victoriametrics.Pusher) map[string]func() error {
	mux.Lock()
	metrics := map[string]func() error{
		"cycletime": func() error { return pusher.ImportCycleTimeMetric(cycletime.NewCycleTimeService(gitlab)) },
		"fridaymergerequests": func() error {
			return pusher.ImporFridayMergeRequests(fridaymergerequests.NewFridayMergeRequests(gitlab))
		},
		"longrunningmergerequests": func() error {
			return pusher.ImportLongTunningMergeRequests(longrunningmergerequests.NewLongRunningMergerequestsService(gitlab))
		},
		"mergerequestcomments": func() error {
			return pusher.ImportMergeRequestComments(mergerequestcomments.NewMergeRequestCommentsService(gitlab))
		},
		"mergerequestparticipants": func() error {
			return pusher.ImportMergeRequestParticipants(mergerequestparticipants.NewMergeRequestParticipantsService(gitlab))
		},
		"mergerequestrate": func() error {
			return pusher.ImportMergeRequestRate(mergerequestrate.NewMergeRequestRateService(gitlab))
		},
		"mergerequestsize": func() error {
			return pusher.ImportMergeRequestSize(mergerequestsize.NewMergeRequestSizeService(gitlab))
		},
		"selfmerging": func() error { return pusher.ImportSelfMerging(selfmerging.NewSelfMergingService(gitlab)) },
		"unreviewedmergerequests": func() error {
			return pusher.ImportUnreviewedMergeRequests(unreviewedmergerequests.NewUnreviewedMergerequests(gitlab))
		},
		"mergerequestthroughput": func() error {
			return pusher.ImportMergeRequestThroughput(mergerequestthroughput.NewMergeRequestThroughputService(gitlab))
		},
		"defectrate": func() error {
			return pusher.ImportDefectRate(defectrate.NewDefectRateService(gitlab))
		},
		"userdefectrate": func() error {
			return pusher.ImportUserDefectRate(defectrate.NewDefectRateService(gitlab))
		},
		"mergerequestsuccessrate": func() error {
			return pusher.ImportMergeRequestSuccessRate(mergerequestsuccessrate.NewMergeRequestSuccessRateService(gitlab))
		},
		"activecontributors": func() error {
			return pusher.ImportActiveContributors(activecontributors.NewActiveContributors(gitlab))
		},
		"mergerequestengagement": func() error {
			return pusher.ImportMergeRequestEngagement(mergerequestengagement.NewMergeRequestEngagementService(gitlab))
		},
	}
	mux.Unlock()
	return metrics
}

func executeMetricFunctions(metrics map[string]func() error) {
	log.Info().Msg("begin ingestion")
	var wg sync.WaitGroup
	mux.RLock()
	for name, f := range metrics {
		wg.Add(1)
		go func(wg *sync.WaitGroup, name string, fn func() error) {
			defer wg.Done()
			log.Info().Str("metric", name).Msg("data ingestion starting")
			err := fn() //call metric fuction
			if err != nil {
				log.Err(err).Str("metric", name).Msg("metric does not write")
			} else {
				log.Info().Str("metric", name).Msg("metric write successfully")
			}
		}(&wg, name, f)
	}
	mux.RUnlock()
	wg.Wait()
	log.Info().Msg("end ingestion")
}

func run(interval time.Duration, gitlab *gitlab.SCM, pusher *victoriametrics.Pusher) {
	metricfunctions := setMetricsFunctions(mux, gitlab, pusher)
	ticker := time.NewTicker(interval)
	for range ticker.C {
		ticker.Stop()
		executeMetricFunctions(metricfunctions)
		ticker.Reset(interval)
	}
}
