package works

import "swe-dashboard/internal/models"

type SCM interface {
	ListProjects() (repositories []*models.Repo, err error)
	ListCommits(projectID int, createdafterday int) (commits []*models.Commit, err error)
}

type WorksService interface {
	NewWorks() (counts []models.ItemCount, err error)
	LegacyRefactor() []models.ItemCount
	HelpOthers() []models.ItemCount
	Churn() []models.ItemCount
}

func NewWorksService(scm SCM) WorksService {
	return &worksService{scm: scm}
}

type worksService struct {
	scm SCM
}

func (w *worksService) NewWorks() (counts []models.ItemCount, err error) {
	return counts, err
}

func (w *worksService) LegacyRefactor() (counts []models.ItemCount) {
	return counts
}

func (w *worksService) HelpOthers() (counts []models.ItemCount) {
	return counts
}

func (w *worksService) Churn() (counts []models.ItemCount) {
	return counts
}
