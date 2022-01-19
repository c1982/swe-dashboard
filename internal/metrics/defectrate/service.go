package defectrate

import (
	"strings"
	"swe-dashboard/internal/models"
)

type SCM interface {
	ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error)
	GetRepository(projectID int) (repository models.Repo, err error)
}

type DefectRateService interface {
}

type defectrate struct {
	scm    SCM
	labels []string
}

func (d *defectrate) NewDefectRateService(scm SCM) DefectRateService {
	return &defectrate{scm: scm,
		labels: []string{"fix", "revert", "bug", "bugfix", "repair", "refactor"}}
}

func (d *defectrate) List() (defects []models.ItemCount, err error) {
	mergerequests, err := d.scm.ListMergeRequest("merged", "all", 1)
	if err != nil {
		return defects, err
	}

	defects = []models.ItemCount{}
	repositories := mergerequests.GroupByRepositories()
	for i := 0; i < len(repositories); i++ {
		repo, err := d.scm.GetRepository(repositories[i].ID)
		if err != nil {
			return defects, err
		}

		repo.MRs = repositories[i].MRs
		defectcount := 0
		for n := 0; n < len(repo.MRs); n++ {
			mr := repo.MRs[n]
			isdefect := d.isDefectMergeRequest(mr.Title)
			if !isdefect {
				continue
			}
			defectcount++
		}

		defects = append(defects, models.ItemCount{
			Name:  repo.Name,
			Count: float64(defectcount),
		})
	}

	return defects, nil
}

func (d *defectrate) isDefectMergeRequest(title string) bool {
	for i := 0; i < len(d.labels); i++ {
		suffix := d.labels[i]
		if strings.HasSuffix(title, suffix) {
			return true
		}
	}
	return false
}
