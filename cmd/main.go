package main

import (
	"fmt"
	"os"
	"swe-dashboard/internal/metrics/mergerequestsize"
	"swe-dashboard/internal/scm/gitlab"
)

func main() {
	baseURL := os.Getenv("SWE_DASHBOARD_GITLAB_BASEURL")
	token := os.Getenv("SWE_DASHBOARD_GITLAB_TOKEN")

	gitlab, err := gitlab.NewSCM(gitlab.GitlabBaseURL(baseURL), gitlab.GitlabToken(token))
	if err != nil {
		panic(err)
	}

	mrsizes := mergerequestsize.NewMergeRequestSizeService(gitlab)
	sizes, err := mrsizes.MergeRequestSizes("merged", "all", 10)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(sizes); i++ {
		fmt.Printf("%s\t%f\r\n", sizes[i].Date, sizes[i].Count)
	}
}
