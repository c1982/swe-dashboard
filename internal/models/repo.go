package models

type Repo struct {
	ID   int
	Name string
	MRs  MergeRequests
}
