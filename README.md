# swe-dashboard
Software Engineer Metrics

## Configuration

* scm-gitlab-baseurl
* scm-gitlab-token
* victoriametrics-importurl
* check-interval

## Runing

Download binary from [releases](https://github.com/c1982/swe-dashboard/releases)

```bash
./swed --scm-gitlab-baseurl=https://your-domain-name/api/v4 \
--scm-gitlab-token=TOKEN \
--victoriametrics-importurl=http://localhost:8428/api/v1/import/prometheus \
--check-interval=1h
```

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
#### Daemon

If you want to run the daemon as a service, you can use the following commands:

0. edit [.swed.config](./daemon/.swed.config) file for your system
1. mkdir /opt/swed
2. cd /opt/swed
3. copy .swed-config /opt/swed
4. copy swed binary to /opt/swed
5. `chmod +x /opt/swed/swed`
6. copy [.swed.config](./daemon/swed.service) to /etc/systemd/system
7. `systemctl enable swed.service`
8. `systemctl start swed.service`
9. watch for errors `journalctl -u swed.service -f`

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
* Active Contributors
* Additin/Deletion Lines of Code

## Supported SCM

* Gitlab Community Edition
* GitHub (implemented, not tested)

## Supported TimeSeries DB

* Victoriametrics

### Grafana Dashboards

* [SWE Dashboard - Metrics](./grafana/swe-dashboard-metrics.json)
* [SWE Dashboard - Repository](./grafana/swe-dashboard-repository.json)

### Metrics Dashboard:

![](docs/images/main-dashboard.png)

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

### Active Contributors

![](docs/images/active-contributors.png)
