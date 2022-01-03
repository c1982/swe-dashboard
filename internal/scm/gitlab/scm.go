package gitlab

import (
	"swe-dashboard/internal/models"
	"time"

	"github.com/xanzy/go-gitlab"
)

const (
	mrStateMerged    = "merged"
	perPageItemCount = 25
)

type GitlabOption func(*SCM) error

func GitlabToken(token string) GitlabOption {
	return func(s *SCM) error {
		return s.setToken(token)
	}
}

func GitlabBaseURL(baseuri string) GitlabOption {
	return func(s *SCM) error {
		return s.setBaseURL(baseuri)
	}
}

type SCM struct {
	client  *gitlab.Client
	token   string
	baseURL string
}

func NewSCM(options ...GitlabOption) (scm *SCM, err error) {
	scm = &SCM{}
	for _, fn := range options {
		if fn == nil {
			continue
		}
		if err := fn(scm); err != nil {
			return scm, err
		}
	}

	c, err := gitlab.NewClient(scm.token, gitlab.WithBaseURL(scm.baseURL))
	if err != nil {
		return scm, err
	}
	scm.client = c

	return scm, nil
}

func (s *SCM) ListRepositories() {

}

func (s *SCM) ListUsers() (users models.Users, err error) {
	opt := &gitlab.ListUsersOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: perPageItemCount,
			Page:    1,
		},
	}

	users = models.Users{}
	for {
		userlist, resp, err := s.client.Users.ListUsers(opt)
		if err != nil {
			return users, err
		}

		for i := 0; i < len(userlist); i++ {
			u := userlist[i]
			activity := time.Now()
			if u.LastActivityOn != nil {
				activity = time.Time(*u.LastActivityOn)
			}

			user := models.User{
				ID:             u.ID,
				Username:       u.Username,
				Name:           u.Name,
				Email:          u.Email,
				State:          u.State,
				AvatarURL:      u.AvatarURL,
				IsAdmin:        u.IsAdmin,
				CreatedAt:      *u.CreatedAt,
				LastActivityOn: activity,
			}
			users[u.ID] = user
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return users, nil
}

func (s *SCM) ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error) {
	mergerequests = []models.MergeRequest{}
	createafter := -1 * ((time.Hour * 24) * time.Duration(createdafterday))
	opt := &gitlab.ListMergeRequestsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: perPageItemCount,
			Page:    1,
		},
		CreatedAfter: gitlab.Time(time.Now().Add(createafter)),
	}

	if scope != "" {
		opt.Scope = gitlab.String(scope)
	}

	if state != "" {
		opt.State = gitlab.String(state)
	}

	for {
		list, rsp, err := s.client.MergeRequests.ListMergeRequests(opt)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(list); i++ {
			mr := list[i]
			asignees := s.convertBasicUsersToUsers(mr.Assignees)
			reviewers := s.convertBasicUsersToUsers(mr.Reviewers)
			mergerequests = append(mergerequests, models.MergeRequest{
				ID:           mr.ID,
				IID:          mr.IID,
				TargetBranch: mr.TargetBranch,
				SourceBranch: mr.SourceBranch,
				ProjectID:    mr.ProjectID,
				Title:        mr.Title,
				State:        mr.State,
				CreatedAt:    *mr.CreatedAt,
				UpdatedAt:    *mr.UpdatedAt,
				Assignee:     s.convertBasicUserToUser(mr.Assignee),
				MergedBy:     s.convertBasicUserToUser(mr.MergedBy),
				ClosedBy:     s.convertBasicUserToUser(mr.ClosedBy),
				MergedAt:     mr.MergedAt,
				ClosedAt:     mr.ClosedAt,
				Assignees:    asignees,
				Reviewers:    reviewers,
			})
		}

		if rsp.NextPage == 0 {
			break
		}

		opt.Page = rsp.NextPage
	}

	return mergerequests, nil
}

func (s *SCM) GetMergeRequestChanges(projectID int, mergeRequestID int) (mergerequest models.MergeRequest, err error) {

	opts := &gitlab.GetMergeRequestChangesOptions{
		AccessRawDiffs: gitlab.Bool(true),
	}

	mr, _, err := s.client.MergeRequests.GetMergeRequestChanges(projectID, mergeRequestID, opts)
	if err != nil {
		return mergerequest, err
	}

	assignees := s.convertBasicUsersToUsers(mr.Assignees)
	reviewers := s.convertBasicUsersToUsers(mr.Reviewers)
	changes := s.convertMergeRequestChanges(mr)
	mergerequest = models.MergeRequest{
		ID:           mr.ID,
		IID:          mr.IID,
		TargetBranch: mr.TargetBranch,
		SourceBranch: mr.SourceBranch,
		ProjectID:    mr.ProjectID,
		Title:        mr.Title,
		State:        mr.State,
		CreatedAt:    *mr.CreatedAt,
		UpdatedAt:    *mr.UpdatedAt,
		Assignee:     s.convertBasicUserToUser(mr.Assignee),
		Assignees:    assignees,
		Reviewers:    reviewers,
		MergedBy:     s.convertBasicUserToUser(mr.MergedBy),
		MergedAt:     mr.MergedAt,
		ClosedBy:     s.convertBasicUserToUser(mr.ClosedBy),
		ClosedAt:     mr.ClosedAt,
		Changes:      changes,
	}

	return mergerequest, nil
}

func (s *SCM) convertBasicUsersToUsers(basicusers []*gitlab.BasicUser) []*models.User {
	users := []*models.User{}
	for i := 0; i < len(basicusers); i++ {
		u := basicusers[i]
		users = append(users, &models.User{
			ID:       u.ID,
			Username: u.Username,
			Name:     u.Name,
			State:    u.State,
		})
	}
	return users
}

func (s *SCM) convertBasicUserToUser(basicuser *gitlab.BasicUser) *models.User {
	if basicuser == nil {
		return &models.User{}
	}

	return &models.User{
		ID:       basicuser.ID,
		Username: basicuser.Username,
		Name:     basicuser.Name,
		State:    basicuser.State,
	}
}

func (s *SCM) convertMergeRequestChanges(mergerequest *gitlab.MergeRequest) []*models.MergeRequestChanges {
	changes := []*models.MergeRequestChanges{}
	for i := 0; i < len(mergerequest.Changes); i++ {
		c := mergerequest.Changes[i]
		changes = append(changes, &models.MergeRequestChanges{
			OldPath:     c.OldPath,
			NewPath:     c.NewPath,
			AMode:       c.AMode,
			BMode:       c.BMode,
			Diff:        c.Diff,
			NewFile:     c.NewFile,
			RenamedFile: c.RenamedFile,
			DeletedFile: c.DeletedFile,
		})
	}

	return changes
}

func (s *SCM) setToken(token string) error {
	s.token = token
	return nil
}

func (s *SCM) setBaseURL(baseuri string) error {
	s.baseURL = baseuri
	return nil
}
