# swe-dashboard
Software Engineer Metrics

## Configuration

* scm-gitlab-baseurl
* scm-gitlab-token
* victoriametrics-importurl
* check-interval

## Runing

Download binary from releases

> ./swed --scm-gitlab-baseurl=https://your-domain-name/api/v4 --scm-gitlab-token=TOKEN --victoriametrics-importurl=http://localhost:8428/api/v1/import/prometheus --check-interval=1h

#### Docker

```bash
docker run --rm --name=swed c1982/swed \
--scm-gitlab-baseurl=https://your-domain-name/api/v4 \
--scm-gitlab-token=TOKEN \
--victoriametrics-importurl=http://localhost:8428/api/v1/import/prometheus \
--check-interval=1h
```

#### Easy Setup

Note: Firstly enter your gitlab variables in config.env file

```bash
git clone https://github.com/c1982/swe-dashboard.git
cd swe-dashboard/docker
docker-compose --env-file ./config.env up
```

## Metrics

* Cycletime
* Time to Open
* Time to Review
* Time to Approve
* Time to Merge
* Friday MRs/PRs
* Long-Running MRs/PRs
* MRs/PRs Comments LeaderBoard
* MRs/PRs Participants LeaderBoard
* MRs/PRs Rates
* MRs/PRs Sizes
* MRs/PRs Throughput
* Self-Merging MRs/PRs
* Developer Turnover Rate
* Unreviewed MRs/PRs
* Review Coverage
* Defect Rate
* MRs/PRs Success Rate

## Supported SCM

* Gitlab
* GitHub (not implemented yet)

## Supported TimeSeries DB

* Victoriametrics


### MRs/PRs Cycle Times

![](docs/images/merge-request-cycle-times.png)

### MRs/PRs Times
![](docs/images/merge-request-times.png)

### Single Cycle Times
![](docs/images/merge-request-single-times.png)

### Long-Running MRs/PRs

![](docs/images/long-running-merge-requests.png)

### MRs/PRs Size Counts

![](docs/images/merge-request-size-counts.png)

### Friday MRs/PRs

![](docs/images/friday-merge-requests.png)

### Unreviewed MRs/PRs

![](docs/images/unreviewed-merge-requests.png)

### MRs/PRs Comments

![](docs/images/user-merge-request-comments.png)

### MRs/PRs Rates

![](docs/images/merge-request-rates.png)

### MRs/PRs Participants

![](docs/images/merge-request-participants.png)

### Self-Merging Users

![](docs/images/self-merging-users.png)

### Repositories Cycle Time

![](docs/images/repo-cycle-time.png)

### Repositories Time to Open

![](docs/images/repo-time-to-open.png)

### Repositories Time to Review

![](docs/images/repo-time-to-review.png)

### Repositories Time to Approve

![](docs/images/repo-time-to-approve.png)

### Repositories Time to Merge

![](docs/images/repo-time-to-merge.png)

### Defect Rate

![](docs/images/defect-rate.png)

### User Defect Rate

TODO: image

### MRs/PRs Success Rate

TODO: image