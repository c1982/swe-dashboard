package assetworkingtime

import (
	"fmt"
	"sort"
	"strings"
	"swe-dashboard/internal/models"
	"time"
)

type SCM interface {
	ListProjects() (repositories []*models.Repo, err error)
	ListCommits(projectID int, createdafterday int) (commits []*models.Commit, err error)
	CommitChanges(projectID int, commitID string) (changes []*models.Change, err error)
}

type AssetWorkingTimeService interface {
	CalculateChanges() (err error)
	Weights() []models.ItemCount
	Iterations() []models.ItemCount
	WorkingHours() []models.ItemCount
}

type change struct {
	name            string
	project         string
	weight          int
	iterations      int
	workingtimeHour int64
	commitTime      time.Time
}

func NewAssetWorkingTimeService(scm SCM, assetExtensions ...string) AssetWorkingTimeService {
	srv := &workingTimeService{
		scm:        scm,
		extensions: assetExtensions,
		changes:    map[string]change{},
	}
	return srv
}

type workingTimeService struct {
	scm        SCM
	extensions []string
	changes    map[string]change
}

func (w *workingTimeService) isAsset(filename string) bool {
	for i := 0; i < len(w.extensions); i++ {
		suffix := w.extensions[i]
		if strings.HasSuffix(filename, suffix) {
			return true
		}
	}
	return false
}

func (w *workingTimeService) Iterations() (iterations []models.ItemCount) {
	iterations = []models.ItemCount{}
	for _, c := range w.changes {
		iterations = append(iterations, models.ItemCount{
			Name:  c.project,
			Name1: c.name,
			Count: float64(c.iterations),
		})
	}
	return iterations
}

func (w *workingTimeService) Weights() (weights []models.ItemCount) {
	weights = []models.ItemCount{}
	for _, c := range w.changes {
		weights = append(weights, models.ItemCount{
			Name:  c.project,
			Name1: c.name,
			Count: float64(c.weight),
		})
	}
	return weights
}

func (w *workingTimeService) WorkingHours() (workinghourse []models.ItemCount) {
	workinghourse = []models.ItemCount{}
	for _, c := range w.changes {
		workinghourse = append(workinghourse, models.ItemCount{
			Name:  c.project,
			Name1: c.name,
			Count: float64(c.workingtimeHour),
		})
	}
	return workinghourse
}

func (w *workingTimeService) CalculateChanges() (err error) {
	projects, err := w.scm.ListProjects()
	if err != nil {
		return err
	}

	list := map[string]change{}
	for _, p := range projects {
		commits, err := w.scm.ListCommits(p.ID, 30)
		if err != nil {
			return err
		}

		sort.Slice(commits, func(i, j int) bool {
			return commits[i].CreatedAt.Before(commits[j].CreatedAt)
		})

		for _, cm := range commits {
			commitchanges, err := w.scm.CommitChanges(p.ID, cm.ShortID)
			if err != nil {
				return err
			}
			for i := 0; i < len(commitchanges); i++ {
				c := commitchanges[i]
				if !w.isAsset(c.Name) {
					continue
				}
				name := fmt.Sprintf("%s-%s", p.Name, c.Name)
				v, ok := list[name]
				if !ok {
					list[name] = change{
						name:            c.Name,
						project:         p.Name,
						weight:          c.Weight,
						iterations:      1,
						workingtimeHour: 0,
						commitTime:      cm.CreatedAt,
					}
				} else {
					v.iterations = v.iterations + 1
					v.weight += c.Weight
					if cm.CreatedAt.After(v.commitTime) {
						v.workingtimeHour = int64(cm.CreatedAt.Sub(v.commitTime) / time.Hour)
					}
					list[name] = v
				}
			}
		}
	}
	w.changes = list
	return nil
}
