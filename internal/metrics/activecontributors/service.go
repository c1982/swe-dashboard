package activecontributors

import "swe-dashboard/internal/models"

type usercount struct {
	Name      string
	Commits   int
	Additions int
	Deletions int
}

type SCM interface {
	ListProjects() (repositories []*models.Repo, err error)
	ListCommits(projectID int, createdafterday int) (commits []*models.Commit, err error)
	ListAllProjectMembers(projectID int) (members []models.User, err error)
}

type ActiveContributorsService interface {
	List() (counts []models.ItemCount, err error)
	Impact() (counts []models.ItemCount)
}

func NewActiveContributors(scm SCM) ActiveContributorsService {
	return &activecontributors{scm: scm}
}

type activecontributors struct {
	scm     SCM
	impacts []models.ItemCount
}

func (a *activecontributors) List() (counts []models.ItemCount, err error) {
	repos, err := a.scm.ListProjects()
	if err != nil {
		return counts, err
	}
	counts = []models.ItemCount{}
	impacts := []models.ItemCount{}
	for i := 0; i < len(repos); i++ {
		r := repos[i]
		commits, err := a.scm.ListCommits(r.ID, 30)
		if err != nil {
			return counts, err
		}
		users, err := a.scm.ListAllProjectMembers(r.ID)
		if err != nil {
			return counts, err
		}
		usercounts := a.commitsCountByUser(commits, users)
		for k, v := range usercounts {
			counts = append(counts, models.ItemCount{
				Name:  r.Name,
				Name1: k,
				Count: float64(v.Commits),
			})
			impacts = append(impacts, models.ItemCount{
				Name:   r.Name,
				Name1:  k,
				Count:  float64(v.Additions),
				Count1: float64(v.Deletions),
			})
		}
	}
	a.impacts = impacts
	return counts, err
}

func (a *activecontributors) Impact() (counts []models.ItemCount) {
	return a.impacts
}

func (a *activecontributors) commitsCountByUser(commits []*models.Commit, users []models.User) (counts map[string]*usercount) {
	counts = map[string]*usercount{}
	for i := 0; i < len(commits); i++ {
		c := commits[i]
		ok, user := a.findProjectMember(c.AuthorName, users)
		if !ok {
			continue
		}
		v, ok := counts[user.Name]
		if !ok {
			counts[user.Name] = &usercount{
				Name:      user.Name,
				Commits:   1,
				Additions: c.Additions,
				Deletions: c.Deletions,
			}
		} else {
			v.Commits = v.Commits + 1
			v.Additions = v.Additions + c.Additions
			v.Deletions = v.Deletions + c.Deletions
			counts[user.Name] = v
		}
	}
	return counts
}

func (a *activecontributors) findProjectMember(name string, users []models.User) (ok bool, user models.User) {
	for i := 0; i < len(users); i++ {
		u := users[i]
		if u.Name == name {
			return true, u
		}
	}
	return false, models.User{}
}
