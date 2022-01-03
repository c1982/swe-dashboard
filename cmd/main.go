package main

import (
	"fmt"
	"os"
	"swe-dashboard/internal/metrics/mergerequestrates"
	"swe-dashboard/internal/scm/gitlab"
)

func main() {
	baseURL := os.Getenv("SWE_DASHBOARD_GITLAB_BASEURL")
	token := os.Getenv("SWE_DASHBOARD_GITLAB_TOKEN")

	gitlab, err := gitlab.NewSCM(gitlab.GitlabBaseURL(baseURL), gitlab.GitlabToken(token))
	if err != nil {
		panic(err)
	}

	mrrates := mergerequestrates.NewMergeRequestRateService(gitlab)
	rates, err := mrrates.MergeRequestRates("merged", "all", 10)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(rates); i++ {
		fmt.Printf("%s\t%f\r\n", rates[i].Date, rates[i].Count)
	}

}
