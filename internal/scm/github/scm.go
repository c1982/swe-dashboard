package github

import (
	"context"
	"swe-dashboard/internal/models"
	"time"

	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

type SCM struct {
	client        *github.Client
	token         string
	baseURL       string
	ctx           context.Context
	organizations []string
}

const (
	perPageItemCount = 25
)

func NewSCM(options ...GithubOption) (scm *SCM, err error) {
	scm = &SCM{}
	for _, fn := range options {
		if fn == nil {
			continue
		}
		if err := fn(scm); err != nil {
			return scm, err
		}
	}

	ctx := context.Background()
	scm.ctx = ctx
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: scm.token},
	)
	tc := oauth2.NewClient(ctx, ts)
	c := github.NewClient(tc)

	scm.client = c

	return scm, nil
}

func (s *SCM) GetSelfOrganizations() (orgs []*github.Organization, err error) {
	opt := &github.ListOptions{
		PerPage: perPageItemCount,
	}

	organizations := []*github.Organization{}

	for {
		data, rsp, err := s.client.Organizations.List(s.ctx, "", opt)
		if err != nil {
			return organizations, err
		}

		organizations = append(organizations, data...)

		if rsp.NextPage == 0 {
			break
		}

		opt.Page = rsp.NextPage
	}

	return organizations, nil
}

func (s *SCM) listOrganizations() (orgs []*github.Organization, err error) {
	if len(s.organizations) > 0 {
		organizations := []*github.Organization{}
		for _, orgName := range s.organizations {
			organization, _, err := s.client.Organizations.Get(s.ctx, orgName)
			if err != nil {
				return organizations, err
			}

			organizations = append(organizations, organization)
		}

		return organizations, nil
	}
	return s.GetSelfOrganizations()
}

func (s *SCM) OrganizationRepositoriesList(orgName string) (repos []*github.Repository, err error) {
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{
			PerPage: perPageItemCount,
			Page:    1,
		},
	}

	repos = []*github.Repository{}

	for {
		data, rsp, err := s.client.Repositories.ListByOrg(s.ctx, orgName, opt)
		if err != nil {
			return repos, err
		}

		repos = append(repos, data...)

		if rsp.NextPage == 0 {
			break
		}

		opt.Page = rsp.NextPage
	}

	return repos, nil
}

func (s *SCM) GetRepository(projectID int) (repository models.Repo, err error) {

	repo, _, err := s.client.Repositories.GetByID(s.ctx, int64(projectID))
	if err != nil {
		return repository, err
	}

	repository = models.Repo{
		ID:             projectID,
		Name:           repo.GetName(),
		Description:    repo.GetDescription(),
		CreatorID:      int(repo.GetOwner().GetID()),
		LastActivityAt: &repo.UpdatedAt.Time,
	}

	return repository, nil
}

func (s *SCM) ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error) {
	mergerequests = []models.MergeRequest{}
	filterdate := time.Now().AddDate(0, 0, -createdafterday)
	mergedOnly := false
	if state == "merged" {
		mergedOnly = true
		state = "closed"
	}
	opt := &github.PullRequestListOptions{
		State:     state,
		Sort:      "created",
		Direction: "desc",
		ListOptions: github.ListOptions{
			PerPage: perPageItemCount,
		},
	}
	organizations, err := s.listOrganizations()
	if err != nil {
		return mergerequests, err
	}

	for _, org := range organizations {
		orgLogin := org.GetLogin()
		repos, err := s.OrganizationRepositoriesList(orgLogin)
		if err != nil {
			return mergerequests, err
		}

		for _, repo := range repos {
			for {
				list, rsp, err := s.client.PullRequests.List(s.ctx, orgLogin, repo.GetName(), opt)
				if err != nil {
					return mergerequests, err
				}

				if len(list) == 0 {
					break
				}

				for _, v := range list {
					if v.GetCreatedAt().Unix() < filterdate.Unix() {
						break
					}
					if mergedOnly {
						if v.MergedAt == nil {
							continue
						}
						mergerequests = append(mergerequests, s.convertGithubPullRequestToMergeRequest(v))
					} else {
						mergerequests = append(mergerequests, s.convertGithubPullRequestToMergeRequest(v))
					}
				}

				if rsp.NextPage == 0 {
					break
				}

				opt.Page = rsp.NextPage
			}
		}

	}
	return mergerequests, nil
}

func (s *SCM) ListMergeRequestNotes(projectID int, mergeRequestID int) (comments []*models.Comment, err error) {
	comments = []*models.Comment{}

	opt := &github.ListOptions{
		PerPage: perPageItemCount,
	}

	repo, _, err := s.client.Repositories.GetByID(s.ctx, int64(projectID))

	if err != nil {
		return comments, err
	}

	owner := repo.GetOwner().GetLogin()
	if repo.GetFork() {
		owner = repo.GetParent().GetOwner().GetLogin()
	}

	for {
		data, rsp, err := s.client.PullRequests.ListReviews(s.ctx, owner, repo.GetName(), int(mergeRequestID), opt)
		if err != nil {
			return comments, err
		}

		if len(data) == 0 {
			break
		}

		comments = append(comments, s.convertGithubPullRequestsCommentsToComments(data)...)

		if rsp.NextPage == 0 {
			break
		}

		opt.Page = rsp.NextPage
	}
	return comments, nil
}

func (s *SCM) ListUsers() (users models.Users, err error) {
	users = models.Users{}

	opt := &github.ListMembersOptions{
		ListOptions: github.ListOptions{
			PerPage: perPageItemCount,
		},
	}

	organizations, err := s.listOrganizations()
	if err != nil {
		return users, err
	}

	for _, v := range organizations {
		for {
			data, rsp, err := s.client.Organizations.ListMembers(s.ctx, v.GetLogin(), opt)
			if err != nil {
				return users, err
			}

			for _, v := range s.convertGithubUsersToUsers(data) {
				users[v.ID] = *v
			}

			if rsp.NextPage == 0 {
				break
			}

			opt.Page = rsp.NextPage
		}
	}

	return users, nil
}

func (s *SCM) GetMergeRequestChanges(projectID int, mergeRequestID int) (mergerequest models.MergeRequest, err error) {
	mergerequest = models.MergeRequest{}

	repo, _, err := s.client.Repositories.GetByID(s.ctx, int64(projectID))
	if err != nil {
		return mergerequest, err
	}

	pullRequest, _, err := s.client.PullRequests.Get(s.ctx, repo.Owner.GetLogin(), repo.GetName(), mergeRequestID)
	if err != nil {
		return mergerequest, err
	}

	mergerequest = s.convertGithubPullRequestToMergeRequest(pullRequest)

	opt := &github.ListOptions{
		Page:    1,
		PerPage: perPageItemCount,
	}

	changes := []*models.MergeRequestChanges{}

	for {
		data, rsp, err := s.client.PullRequests.ListFiles(s.ctx, repo.Owner.GetLogin(), repo.GetName(), mergeRequestID, opt)
		if err != nil {
			return mergerequest, err
		}

		changes = append(changes, s.convertGithubPullRequestChangesToMergeRequestChanges(data)...)

		if rsp.NextPage == 0 {
			break
		}

		opt.Page = rsp.NextPage
	}

	mergerequest.Changes = changes

	return mergerequest, nil
}

func (s *SCM) ListAllProjectMembers(projectID int) (members []models.User, err error) {
	members = []models.User{}

	opt := &github.ListContributorsOptions{
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: perPageItemCount,
		},
	}

	repo, _, err := s.client.Repositories.GetByID(s.ctx, int64(projectID))
	if err != nil {
		return members, err
	}

	for {
		data, rsp, err := s.client.Repositories.ListContributors(s.ctx, repo.Owner.GetLogin(), repo.GetName(), opt)
		if err != nil {
			return members, err
		}

		for _, v := range data {
			members = append(members, *s.convertGithubContributorToUser(v))
		}

		if rsp.NextPage == 0 {
			break
		}

		opt.Page = rsp.NextPage
	}

	return members, nil
}

func (s *SCM) GetMergeRequestParticipants(projectID int, mergeRequestID int) (users []*models.User, err error) {
	users = []*models.User{}

	opt := &github.ListOptions{
		Page:    1,
		PerPage: perPageItemCount,
	}

	repo, _, err := s.client.Repositories.GetByID(s.ctx, int64(projectID))
	if err != nil {
		return users, err
	}

	for {
		reviews, rsp, err := s.client.PullRequests.ListReviews(s.ctx, repo.Owner.GetLogin(), repo.GetName(), mergeRequestID, opt)
		if err != nil {
			return users, err
		}

		for _, v := range reviews {
			users = append(users, s.convertGithubUserToUser(v.GetUser()))
		}

		if rsp.NextPage == 0 {
			break
		}

		opt.Page = rsp.NextPage
	}

	return users, nil
}

func (s *SCM) GetMergeRequestCommits(projectID, mergeRequestID int) (commits []*models.Commit, err error) {
	commits = []*models.Commit{}
	repo, _, err := s.client.Repositories.GetByID(s.ctx, int64(projectID))
	if err != nil {
		return commits, err
	}

	opt := &github.ListOptions{
		Page:    1,
		PerPage: perPageItemCount,
	}

	for {
		data, rsp, err := s.client.PullRequests.ListCommits(s.ctx, repo.Owner.GetLogin(), repo.GetName(), mergeRequestID, opt)
		if err != nil {
			return commits, err
		}

		for _, v := range data {
			commit := s.convertGithubCommitToCommit(v, projectID)
			commits = append(commits, commit)
		}

		if rsp.NextPage == 0 {
			break
		}

		opt.Page = rsp.NextPage
	}

	return commits, nil
}

func (s *SCM) convertGithubCommitToCommit(repoCommit *github.RepositoryCommit, repoID int) *models.Commit {
	commit := repoCommit.GetCommit()
	return &models.Commit{
		ID:             commit.GetNodeID(),
		Message:        commit.GetMessage(),
		AuthorName:     repoCommit.GetAuthor().GetName(),
		CommitterName:  commit.GetAuthor().GetName(),
		CommitterEmail: repoCommit.GetAuthor().GetEmail(),
		CommittedDate:  commit.GetCommitter().GetDate(),
		CreatedAt:      commit.GetAuthor().GetDate(),
		ProjectID:      repoID,
	}
}

func (s *SCM) convertGithubPullRequestsCommentsToComments(comments []*github.PullRequestReview) []*models.Comment {
	var commentsList []*models.Comment

	for _, comment := range comments {

		//TODO: need to add more fields
		commentsList = append(commentsList, &models.Comment{
			ID:        int(comment.GetID()),
			Body:      comment.GetBody(),
			Author:    *s.convertGithubUserToUser(comment.User),
			CreatedAt: comment.GetSubmittedAt(),
		})
	}

	return commentsList
}

func (s *SCM) convertGithubUsersToUsers(users []*github.User) []*models.User {
	usersList := []*models.User{}

	for _, user := range users {
		usersList = append(usersList, s.convertGithubUserToUser(user))
	}

	return usersList
}

func (s *SCM) convertGithubUserToUser(user *github.User) *models.User {
	return &models.User{
		ID:             int(user.GetID()),
		Name:           user.GetName(),
		Username:       user.GetLogin(),
		Email:          user.GetEmail(),
		IsAdmin:        user.GetSiteAdmin(),
		AvatarURL:      user.GetAvatarURL(),
		CreatedAt:      user.GetCreatedAt().Time,
		LastSignInAt:   user.GetUpdatedAt().Time,
		LastActivityOn: user.GetUpdatedAt().Time,
		State:          "",
	}
}

func (s *SCM) convertGithubPullRequestToMergeRequest(pr *github.PullRequest) models.MergeRequest {

	state := pr.GetState()
	//TODO: Check this trick for github merged state
	if pr.MergedAt != nil {
		state = "merged"
	}

	return models.MergeRequest{
		ID:           int(pr.GetID()),
		IID:          int(pr.GetNumber()),
		TargetBranch: pr.GetHead().GetRef(),
		SourceBranch: pr.GetBase().GetRef(),
		ProjectID:    int(pr.GetBase().GetRepo().GetID()),
		Title:        pr.GetTitle(),
		State:        state,
		CreatedAt:    pr.GetCreatedAt(),
		UpdatedAt:    pr.GetUpdatedAt(),
		Assignee:     s.convertGithubUserToUser(pr.Assignee),
		MergedBy:     s.convertGithubUserToUser(pr.MergedBy),
		//TODO: Check this i am not sure about that closed by is correct value.
		ClosedBy:  s.convertGithubUserToUser(pr.MergedBy),
		MergedAt:  pr.MergedAt,
		ClosedAt:  pr.ClosedAt,
		Assignees: s.convertGithubUsersToUsers(pr.Assignees),
		Reviewers: s.convertGithubUsersToUsers(pr.RequestedReviewers),
		Author:    s.convertGithubUserToUser(pr.GetUser()),
		Draft:     pr.GetDraft(),
	}
}

func (s *SCM) convertGithubPullRequestsToMergeRequests(prs []*github.PullRequest) []models.MergeRequest {
	mergeRequests := []models.MergeRequest{}

	for _, pr := range prs {
		mergeRequests = append(mergeRequests, s.convertGithubPullRequestToMergeRequest(pr))
	}

	return mergeRequests
}

func (s *SCM) convertGithubPullRequestChangesToMergeRequestChanges(changes []*github.CommitFile) []*models.MergeRequestChanges {
	mergeRequestChanges := []*models.MergeRequestChanges{}

	for _, change := range changes {
		isFileNameChanged := false

		//TODO: need to check this i am not sure about this is correct way to do it.
		if change.GetFilename() != change.GetPreviousFilename() {
			isFileNameChanged = true
		}

		mergeRequestChanges = append(mergeRequestChanges, &models.MergeRequestChanges{
			Diff:        change.GetPatch(),
			NewPath:     change.GetFilename(),
			OldPath:     change.GetPreviousFilename(),
			RenamedFile: isFileNameChanged,
			//TODO: maybe we can check status for detect if file is added, deleted or something.
		})
	}

	return mergeRequestChanges
}

func (s *SCM) convertGithubContributorToUser(contributor *github.Contributor) *models.User {
	return &models.User{
		ID:        int(*contributor.ID),
		Name:      contributor.GetName(),
		Username:  contributor.GetLogin(),
		Email:     contributor.GetEmail(),
		IsAdmin:   contributor.GetSiteAdmin(),
		AvatarURL: contributor.GetAvatarURL(),
		State:     "",
	}
}

func (s *SCM) setToken(token string) error {
	s.token = token
	return nil
}

func (s *SCM) setBaseURL(baseuri string) error {
	s.baseURL = baseuri
	return nil
}

func (s *SCM) setOrganizations(organizations []string) error {
	s.organizations = organizations
	return nil
}
