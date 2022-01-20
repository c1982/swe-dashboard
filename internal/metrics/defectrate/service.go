package defectrate

import (
	"strings"
	"swe-dashboard/internal/models"
	"time"
)

var (
	defectPrefixes = []string{"fix", "revert", "bug", "bugfix", "repair", "refactor"}
)

type SCM interface {
	ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error)
	GetRepository(projectID int) (repository models.Repo, err error)
}

type DefectRateService interface {
	List() (defects []models.ItemCount, err error)
}

type defectrate struct {
	scm    SCM
	labels []string
}

func NewDefectRateService(scm SCM) DefectRateService {
	return &defectrate{scm: scm,
		labels: defectPrefixes}
}

func (d *defectrate) List() (defects []models.ItemCount, err error) {
	mergerequests, err := d.scm.ListMergeRequest("merged", "all", time.Now().Day())
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

		rate := defectcount * 100 / len(repositories[i].MRs)
		defects = append(defects, models.ItemCount{
			Name:  repo.Name,
			Count: float64(rate),
		})
	}

	return defects, nil
}

func (d *defectrate) isDefectMergeRequest(title string) bool {
	title = strings.ToLower(title)
	for i := 0; i < len(d.labels); i++ {
		suffix := d.labels[i]
		if strings.HasPrefix(title, suffix) {
			return true
		}
	}
	return false
}
