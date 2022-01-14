# swe-dashboard
Software Engineer Metrics

## Configuration

* scm-gitlab-baseurl
* scm-gitlab-token
* victoriametrics-importurl
* check-interval

## Docker

### Build

> docker build -t swed .

### Run

> docker run --rm swed

## Metrics

* Cycletime
* TimeToOpen
* TimeToReview
* TimeToApprove
* TimeToMerge
* Friday MR/PR
* Long-Running MR/PR
* MR/PR Comments LeaderBoard
* MR/PR Participants LeaderBoard
* MR/PR Rates
* MR/PR Sizes
* MR/PR Throughput
* Self-Merging MR/PR
* Developer Turnover Rate
* Unreviewed MR/PR

## Supported SCM

* Gitlab

## Supported TimeSeries DB

* Victoriametrics