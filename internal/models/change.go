package models

type Change struct {
	ProjectID int    `json:"project_id"`
	Name      string `json:"name"`
	Weight    int    `json:"weight"`
}
