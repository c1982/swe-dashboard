package victoriametrics

import (
	"fmt"
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
	cycleTimeMetricName               = `cycle_time{repository="%s", title="%s"} %f`
	timeToOpenMetricName              = `time_to_open{repository="%s", title="%s"} %f`
	timetoReviewMetricName            = `time_to_review{repository="%s", title="%s"} %f`
	timetoApproveMetricName           = `time_to_approve{repository="%s", title="%s"} %f`
	timetoMergeMetricName             = `time_to_merge{repository="%s", title="%s"} %f`
	fridayMergeRequestMetricName      = `friday_merge_request{repository="%s"} %f`
	longRunningMergeRequestMetricName = `long_running_merge_request{repository="%s", title="%s"} %f`
	commentsLeaderboardMetricName     = `comments_leaderboard{id="%d",user="%s", name="%s"} %f`
	participantsLeaderboardMetricName = `participants_leaderboard{id="%d",user="%s", name="%s"} %f`
	mergeRequestRateMetricName        = `merge_request_rate{repository="%s"} %f`
	mergeRequestSizeMetricName        = `merge_request_size{repository="%s", title="%s"} %f`
	mergeRequestThroughputMetricName  = `merge_request_throughput{repository="%s"} %f`
	selfMergingMetricName             = `self_merging{name="%s", username="%s"} %f`
	turnOverRateMetricName            = `turn_over_rate{} %f`
	unreviewedMergeRequestMetricName  = `unreviewed_merge_request{repository="%s"} %f`
	defectRateMetricName              = `defect_rate{repository="%s"} %f`
	userDefectRateMetricName          = `defect_rate_user{repository="%s", username="%s", name="%s"} %f`
	mergeRequestSuccessRateMetricName = `merge_request_success_rate{repository="%s"} %f`
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
	for i := 0; i < len(timetoopens); i++ {
		payload := fmt.Sprintf(timetoReviewMetricName, timetoreview[i].Name, timetoreview[i].Name1, timetoreview[i].Count)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}

	timetoapprove := service.TimeToApprove()
	for i := 0; i < len(timetoopens); i++ {
		payload := fmt.Sprintf(timetoApproveMetricName, timetoapprove[i].Name, timetoapprove[i].Name1, timetoapprove[i].Count)
		err := p.Push(payload)
		if err != nil {
			return err
		}
	}

	timetomerge := service.TimeToMerge()
	for i := 0; i < len(timetoopens); i++ {
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
	leaderboard, err := service.ParticipantsLeaderBoard()
	if err != nil {
		return err
	}

	for i := 0; i < len(leaderboard); i++ {
		payload := fmt.Sprintf(participantsLeaderboardMetricName,
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

func (p *Pusher) ImportMergeRequestRate(service mergerequestrate.MergeRequestRateService) (err error) {
	counts, err := service.MergeRequestRatesThisMonth()
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

//TODO: don't use until fixed
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
