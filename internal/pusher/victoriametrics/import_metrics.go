package victoriametrics

import (
	"fmt"
	"swe-dashboard/internal/metrics/activecontributors"
	"swe-dashboard/internal/metrics/assetiterations"
	"swe-dashboard/internal/metrics/cycletime"
	"swe-dashboard/internal/metrics/defectrate"
	"swe-dashboard/internal/metrics/fridaymergerequests"
	"swe-dashboard/internal/metrics/longrunningmergerequests"
	"swe-dashboard/internal/metrics/mergerequestcomments"
	"swe-dashboard/internal/metrics/mergerequestparticipants"
	"swe-dashboard/internal/metrics/mergerequestrate"
	"swe-dashboard/internal/metrics/mergerequestsize"
	"swe-dashboard/internal/metrics/mergerequestsuccessrate"
	"swe-dashboard/internal/metrics/mergerequestthroughput"
	"swe-dashboard/internal/metrics/selfmerging"
	"swe-dashboard/internal/metrics/turnoverrate"
	"swe-dashboard/internal/metrics/unreviewedmergerequests"
)

const (
	cycleTimeMetricName                          = `cycle_time{repository="%s", title="%s"} %f`
	timeToOpenMetricName                         = `time_to_open{repository="%s", title="%s"} %f`
	timetoReviewMetricName                       = `time_to_review{repository="%s", title="%s"} %f`
	timetoApproveMetricName                      = `time_to_approve{repository="%s", title="%s"} %f`
	timetoMergeMetricName                        = `time_to_merge{repository="%s", title="%s"} %f`
	fridayMergeRequestMetricName                 = `friday_merge_request{repository="%s"} %f`
	longRunningMergeRequestMetricName            = `long_running_merge_request{repository="%s", title="%s"} %f`
	commentsLeaderboardMetricName                = `comments_leaderboard{id="%d",user="%s", name="%s"} %f`
	mergeRequestRateMetricName                   = `merge_request_rate{repository="%s"} %f`
	mergeRequestSizeMetricName                   = `merge_request_size{repository="%s", title="%s"} %f`
	mergeRequestThroughputMetricName             = `merge_request_throughput{repository="%s"} %f`
	selfMergingMetricName                        = `self_merging{name="%s", username="%s"} %f`
	turnOverRateMetricName                       = `turn_over_rate{} %f`
	unreviewedMergeRequestMetricName             = `unreviewed_merge_request{repository="%s"} %f`
	defectRateMetricName                         = `defect_rate{repository="%s"} %f`
	userDefectRateMetricName                     = `defect_rate_user{repository="%s", username="%s", name="%s"} %f`
	mergeRequestSuccessRateMetricName            = `merge_request_success_rate{repository="%s"} %f`
	activeContributorsMetricName                 = `active_contributors{repository="%s", author="%s", email="%s"} %f`
	commitAdditionsMetricName                    = `commit_additions{repository="%s", author="%s", email="%s"} %f`
	commitDeletionsMetricName                    = `commit_deletions{repository="%s", author="%s", email="%s"} %f`
	mergeRequestParticipantsdMetricName          = `merge_request_participants{repository="%s",user="%s", name="%s"} %f`
	mergeRequestEngagementsMetricName            = `merge_request_engagement{repository="%s", author="%s", mergedby="%s"} %f`
	mergeRequestEngagementParticipantsMetricName = `merge_request_engage_participants{repository="%s",author="%s", participant="%s"} %f`

	assetIterationWeights    = `assets_weights{repository="%s",name="%s"} %f`
	assetIterationHours      = `assets_iteration_hours{repository="%s",name="%s"} %f`
	assetIterationIterations = `assets_iterations{repository="%s",name="%s"} %f`
)

func (p *Pusher) ImportCycleTimeMetric(service cycletime.CycleTimeService) (err error) {
	metrics, err := service.CycleTime()
	if err != nil {
		return err
	}

	for i := 0; i < len(metrics); i++ {
		payload := fmt.Sprintf(cycleTimeMetricName, metrics[i].Name, metrics[i].Name1, metrics[i].Count)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}

	timetoopens := service.TimeToOpen()
	for i := 0; i < len(timetoopens); i++ {
		payload := fmt.Sprintf(timeToOpenMetricName, timetoopens[i].Name, timetoopens[i].Name1, timetoopens[i].Count)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}

	timetoreview := service.TimeToReview()
	for i := 0; i < len(timetoreview); i++ {
		payload := fmt.Sprintf(timetoReviewMetricName, timetoreview[i].Name, timetoreview[i].Name1, timetoreview[i].Count)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}

	timetoapprove := service.TimeToApprove()
	for i := 0; i < len(timetoapprove); i++ {
		payload := fmt.Sprintf(timetoApproveMetricName, timetoapprove[i].Name, timetoapprove[i].Name1, timetoapprove[i].Count)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}

	timetomerge := service.TimeToMerge()
	for i := 0; i < len(timetomerge); i++ {
		payload := fmt.Sprintf(timetoMergeMetricName, timetomerge[i].Name, timetomerge[i].Name1, timetomerge[i].Count)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Pusher) ImporFridayMergeRequests(service fridaymergerequests.FridayMergerequestsService) (err error) {
	counts, err := service.List()
	if err != nil {
		return err
	}

	for i := 0; i < len(counts); i++ {
		payload := fmt.Sprintf(fridayMergeRequestMetricName, counts[i].Name, counts[i].Count)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Pusher) ImportLongTunningMergeRequests(service longrunningmergerequests.LongRunningMergerequestsService) (err error) {
	counts, err := service.List()
	if err != nil {
		return err
	}

	for i := 0; i < len(counts); i++ {
		payload := fmt.Sprintf(longRunningMergeRequestMetricName, counts[i].Name, counts[i].Name1, counts[i].Count)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Pusher) ImportMergeRequestComments(service mergerequestcomments.MergeRequestCommentsService) (err error) {
	leaderboard, err := service.List()
	if err != nil {
		return err
	}

	for i := 0; i < len(leaderboard); i++ {
		payload := fmt.Sprintf(commentsLeaderboardMetricName,
			leaderboard[i].ID,
			leaderboard[i].Username,
			leaderboard[i].Name,
			leaderboard[i].Count)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Pusher) ImportMergeRequestParticipants(service mergerequestparticipants.MergeRequestParticipantsService) (err error) {
	leaderboard, err := service.List()
	if err != nil {
		return err
	}

	for i := 0; i < len(leaderboard); i++ {
		payload := fmt.Sprintf(mergeRequestParticipantsdMetricName,
			leaderboard[i].Name,
			leaderboard[i].Name1,
			leaderboard[i].Name2,
			leaderboard[i].Count)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}

	engageparticipants := service.EngageParticipants()
	for i := 0; i < len(engageparticipants); i++ {
		payload := fmt.Sprintf(mergeRequestEngagementParticipantsMetricName,
			engageparticipants[i].Name,
			engageparticipants[i].Name1,
			engageparticipants[i].Name2,
			engageparticipants[i].Count)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}

	engagements := service.Engagements()
	for i := 0; i < len(engagements); i++ {
		payload := fmt.Sprintf(mergeRequestEngagementsMetricName,
			engagements[i].Name,
			engagements[i].Name1,
			engagements[i].Name2,
			engagements[i].Count)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Pusher) ImportMergeRequestRate(service mergerequestrate.MergeRequestRateService) (err error) {
	counts, err := service.MergeRequestRates()
	if err != nil {
		return err
	}

	for i := 0; i < len(counts); i++ {
		payload := fmt.Sprintf(mergeRequestRateMetricName, counts[i].Name, counts[i].Count)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Pusher) ImportMergeRequestSize(service mergerequestsize.MergeRequestSizeService) (err error) {
	counts, err := service.Sizes()
	if err != nil {
		return err
	}

	for i := 0; i < len(counts); i++ {
		payload := fmt.Sprintf(mergeRequestSizeMetricName, counts[i].Name, counts[i].Name1, counts[i].Count)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Pusher) ImportMergeRequestThroughput(service mergerequestthroughput.MergeRequestThroughputService) (err error) {
	counts, err := service.List()
	if err != nil {
		return err
	}

	for i := 0; i < len(counts); i++ {
		payload := fmt.Sprintf(mergeRequestThroughputMetricName, counts[i].Name, counts[i].Count)
		err := p.PushWithTime(payload, counts[i].Date)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Pusher) ImportSelfMerging(service selfmerging.SelfMergingService) (err error) {
	users, err := service.GetSelfMergingUsers()
	if err != nil {
		return err
	}

	for i := 0; i < len(users); i++ {
		payload := fmt.Sprintf(selfMergingMetricName, users[i].Name, users[i].Username, users[i].Count)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO: don't use until fixed
func (p *Pusher) ImportTurnOverRate(service turnoverrate.TurnOverrateService) (err error) {
	counts, err := service.TurnOverRate()
	if err != nil {
		return err
	}

	for i := 0; i < len(counts); i++ {
		payload := fmt.Sprintf(turnOverRateMetricName, counts[i].Count)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Pusher) ImportUnreviewedMergeRequests(service unreviewedmergerequests.UnreviewedMergeRequestsService) (err error) {
	counts, err := service.List()
	if err != nil {
		return err
	}
	for i := 0; i < len(counts); i++ {
		payload := fmt.Sprintf(unreviewedMergeRequestMetricName, counts[i].Name, counts[i].Count)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Pusher) ImportDefectRate(service defectrate.DefectRateService) (err error) {
	counts, err := service.List()
	if err != nil {
		return err
	}

	for i := 0; i < len(counts); i++ {
		payload := fmt.Sprintf(defectRateMetricName, counts[i].Name, counts[i].Count)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}

	return nil
}
func (p *Pusher) ImportUserDefectRate(service defectrate.DefectRateService) (err error) {
	counts, err := service.Users()
	if err != nil {
		return err
	}

	for i := 0; i < len(counts); i++ {
		payload := fmt.Sprintf(userDefectRateMetricName, counts[i].Name, counts[i].Name1, counts[i].Name2, counts[i].Count)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Pusher) ImportMergeRequestSuccessRate(service mergerequestsuccessrate.MergeRequestSuccessRateService) (err error) {
	counts, err := service.List()
	if err != nil {
		return err
	}

	for i := 0; i < len(counts); i++ {
		payload := fmt.Sprintf(mergeRequestSuccessRateMetricName, counts[i].Name, counts[i].Count)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Pusher) ImportActiveContributors(service activecontributors.ActiveContributorsService) (err error) {
	metrics, err := service.List()
	if err != nil {
		return err
	}

	for i := 0; i < len(metrics); i++ {
		payload := fmt.Sprintf(activeContributorsMetricName, metrics[i].Name, metrics[i].Name1, metrics[i].Name2, metrics[i].Count)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}

	impact := service.Impact()
	for i := 0; i < len(impact); i++ {
		payload := fmt.Sprintf(commitAdditionsMetricName, impact[i].Name, impact[i].Name1, impact[i].Name2, impact[i].Count)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}
	for i := 0; i < len(impact); i++ {
		payload := fmt.Sprintf(commitDeletionsMetricName, impact[i].Name, impact[i].Name1, impact[i].Name2, impact[i].Count1)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Pusher) ImportAssertIterations(service assetiterations.AssetIterationTimeService) (err error) {
	err = service.CalculateChanges()
	if err != nil {
		return err
	}

	weights := service.Weights()
	for i := 0; i < len(weights); i++ {
		w := weights[i]
		payload := fmt.Sprintf(assetIterationWeights, w.Name, w.Name1, w.Count)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}

	workinghours := service.IterationHours()
	for i := 0; i < len(workinghours); i++ {
		w := workinghours[i]
		payload := fmt.Sprintf(assetIterationHours, w.Name, w.Name1, w.Count)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}

	iterations := service.Iterations()
	for i := 0; i < len(iterations); i++ {
		w := iterations[i]
		payload := fmt.Sprintf(assetIterationIterations, w.Name, w.Name1, w.Count)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}

	return nil
}
