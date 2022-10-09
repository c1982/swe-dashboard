package gitlab

import (
	"strings"
	"swe-dashboard/internal/models"
	"time"
	"unicode/utf8"

	"github.com/xanzy/go-gitlab"
)

const (
	mrStateMerged    = "merged"
	perPageItemCount = 25
)

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

func (s *SCM) ListProjects() (repositories []*models.Repo, err error) {
	repositories = []*models.Repo{}
	opt := &gitlab.ListProjectsOptions{
		Simple:     gitlab.Bool(true),
		Statistics: gitlab.Bool(true),
	}

	for {
		list, resp, err := s.client.Projects.ListProjects(opt)
		if err != nil {
			return repositories, err
		}
		repositories = append(repositories, s.convertProjectsToRepo(list)...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return repositories, nil
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

func (s *SCM) ListCommits(projectID int, createdafterday int) (commits []*models.Commit, err error) {
	createafter := -1 * ((time.Hour * 24) * time.Duration(createdafterday))
	opt := &gitlab.ListCommitsOptions{
		Since:     gitlab.Time(time.Now().Add(createafter)),
		WithStats: gitlab.Bool(true),
	}

	commits = []*models.Commit{}
	for {
		list, rsp, err := s.client.Commits.ListCommits(projectID, opt)
		if err != nil {
			return commits, err
		}
		commits = append(commits, s.convertCommits(list)...)
		if rsp.NextPage == 0 {
			break
		}
		opt.Page = rsp.NextPage
	}

	return commits, nil
}

func (s *SCM) CommitChanges(projectID int, commitID string) (changes []*models.Change, err error) {
	diffs, _, err := s.client.Commits.GetCommitDiff(projectID, commitID, &gitlab.GetCommitDiffOptions{})
	if err != nil {
		return changes, err
	}

	changes = []*models.Change{}
	for i := 0; i < len(diffs); i++ {
		d := diffs[i]
		if d.DeletedFile {
			continue
		}
		changes = append(changes, &models.Change{
			ProjectID: projectID,
			Name:      d.NewPath,
			Weight:    s.calculateCommitDiffWeight(d.Diff),
		})
	}

	return changes, nil
}

func (s *SCM) calculateCommitDiffWeight(diff string) int {
	lines := strings.Split(diff, "\n")
	var linecount int
	var linesize int

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		line = strings.TrimSpace(line)
		size := utf8.RuneCountInString(line)
		if size == 0 {
			continue
		}
		if line == "" {
			continue
		}
		linesize += size
		linecount++
	}

	if linesize == 0 || linecount == 0 {
		return 0
	}

	weight := linesize / linecount
	return weight
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
				Author:       s.convertBasicUserToUser(mr.Author),
				Draft:        mr.WorkInProgress,
			})
		}

		if rsp.NextPage == 0 {
			break
		}

		opt.Page = rsp.NextPage
	}

	return mergerequests, nil
}

func (s *SCM) GetMergeRequestCommits(projectID, mergeRequestID int) (commits []*models.Commit, err error) {
	opt := &gitlab.GetMergeRequestCommitsOptions{}
	commits = []*models.Commit{}

	for {
		list, rsp, err := s.client.MergeRequests.GetMergeRequestCommits(projectID, mergeRequestID, opt)
		if err != nil {
			return commits, err
		}

		commits = append(commits, s.convertCommits(list)...)

		if rsp.NextPage == 0 {
			break
		}

		opt.Page = rsp.NextPage
	}

	return commits, err
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

func (s *SCM) ListAllProjectMembers(projectID int) (members []models.User, err error) {
	opt := &gitlab.ListProjectMembersOptions{
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: perPageItemCount,
		},
	}

	members = []models.User{}
	for {
		list, rsp, err := s.client.ProjectMembers.ListAllProjectMembers(projectID, opt)
		if err != nil {
			return members, err
		}

		for i := 0; i < len(list); i++ {
			member := list[i]
			user := models.User{
				ID:        member.ID,
				Username:  member.Username,
				Name:      member.Name,
				Email:     member.Email,
				State:     member.State,
				CreatedAt: *member.CreatedAt,
			}
			members = append(members, user)
		}

		if rsp.NextPage == 0 {
			break
		}

		opt.Page = rsp.NextPage
	}

	return members, nil
}

func (s *SCM) GetRepository(projectID int) (repository models.Repo, err error) {
	opt := &gitlab.GetProjectOptions{
		Statistics:           gitlab.Bool(false),
		License:              gitlab.Bool(false),
		WithCustomAttributes: gitlab.Bool(false),
	}

	repo, _, err := s.client.Projects.GetProject(projectID, opt)
	if err != nil {
		return repository, err
	}

	repository = models.Repo{
		ID:             projectID,
		Name:           repo.Name,
		Description:    repo.Description,
		CreatorID:      repo.CreatorID,
		LastActivityAt: repo.LastActivityAt,
	}

	return repository, nil
}

func (s *SCM) GetMergeRequestParticipants(projectID int, mergeRequestID int) (users []*models.User, err error) {
	basicusers, _, err := s.client.MergeRequests.GetMergeRequestParticipants(projectID, mergeRequestID)
	if err != nil {
		return users, err
	}

	users = s.convertBasicUsersToUsers(basicusers)
	return users, nil
}

func (s *SCM) ListMergeRequestNotes(projectID int, mergeRequestID int) (comments []*models.Comment, err error) {
	comments = []*models.Comment{}
	opt := &gitlab.ListMergeRequestNotesOptions{}

	for {
		notes, rsp, err := s.client.Notes.ListMergeRequestNotes(projectID, mergeRequestID, opt)
		if err != nil {
			return comments, err
		}

		for i := 0; i < len(notes); i++ {
			n := notes[i]

			approved := false
			if n.System {
				approved = strings.Contains(n.Body, "approved")
			}

			comments = append(comments, &models.Comment{
				ID:           n.ID,
				Body:         n.Body,
				Title:        n.Title,
				System:       n.System,
				Resolvable:   n.Resolvable,
				Resolved:     n.Resolved,
				UpdatedAt:    *n.UpdatedAt,
				CreatedAt:    *n.CreatedAt,
				NoteableType: n.NoteableType,
				ApprovedNote: approved,
				FileName:     n.FileName,
				Author: models.User{
					ID:       n.Author.ID,
					Name:     n.Author.Name,
					Username: n.Author.Username,
					State:    n.Author.State,
				},
				ResolvedBy: models.User{
					ID:       n.ResolvedBy.ID,
					Name:     n.ResolvedBy.Name,
					Username: n.ResolvedBy.Username,
					State:    n.ResolvedBy.State,
				},
			})
		}

		if rsp.NextPage == 0 {
			break
		}

		opt.Page = rsp.NextPage
	}

	return comments, nil
}

func (s *SCM) convertBasicUsersToUsers(basicusers []*gitlab.BasicUser) []*models.User {
	users := []*models.User{}
	for i := 0; i < len(basicusers); i++ {
		u := basicusers[i]
		users = append(users, &models.User{
			ID:       u.ID,
			Name:     u.Name,
			Username: u.Username,
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
		Name:     basicuser.Name,
		Username: basicuser.Username,
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

func (s *SCM) convertCommits(commits []*gitlab.Commit) []*models.Commit {
	list := []*models.Commit{}
	for i := 0; i < len(commits); i++ {
		c := commits[i]
		commit := &models.Commit{
			ID:             c.ID,
			ShortID:        c.ShortID,
			Title:          c.Title,
			AuthorName:     c.AuthorName,
			AuthorEmail:    c.AuthorEmail,
			CommitterName:  c.CommitterName,
			CommitterEmail: c.CommitterEmail,
			CommittedDate:  *c.CommittedDate,
			CreatedAt:      *c.CreatedAt,
			Message:        c.Message,
			ProjectID:      c.ProjectID,
		}

		if c.Stats != nil {
			commit.Additions = c.Stats.Additions
			commit.Deletions = c.Stats.Deletions
			commit.Total = c.Stats.Total
		}

		list = append(list, commit)
	}

	return list
}

func (s *SCM) convertProjectsToRepo(projects []*gitlab.Project) []*models.Repo {
	repos := []*models.Repo{}
	for i := 0; i < len(projects); i++ {
		p := projects[i]
		r := &models.Repo{
			ID:             p.ID,
			Name:           p.Name,
			Description:    p.Description,
			CreatorID:      p.CreatorID,
			LastActivityAt: p.LastActivityAt,
			CreatedAt:      p.CreatedAt,
		}
		if p.Statistics != nil {
			r.CommitCount = p.Statistics.CommitCount
		}

		repos = append(repos, r)
	}
	return repos
}

func (s *SCM) setToken(token string) error {
	s.token = token
	return nil
}

func (s *SCM) setBaseURL(baseuri string) error {
	s.baseURL = baseuri
	return nil
}
