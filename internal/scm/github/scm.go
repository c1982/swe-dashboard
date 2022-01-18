package github

import (
	"context"
	"swe-dashboard/internal/models"

	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

type SCM struct {
	client  *github.Client
	token   string
	baseURL string
	ctx     context.Context
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

func (s *SCM) listOrganizations() (orgs []*github.Organization, err error) {
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
	//TODO: add createdafterday for filtering date.
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
				list, rsp, _ := s.client.PullRequests.List(s.ctx, orgLogin, repo.GetName(), opt)

				if len(list) == 0 {
					break
				}

				mergerequests = append(mergerequests, s.convertGithubPullRequestsToMergeRequests(list)...)

				if rsp.NextPage == 0 {
					break
				}

				opt.ListOptions.Page = rsp.NextPage
			}
		}

	}

	return mergerequests, nil
}

func (s *SCM) ListMergeRequestNotes(projectID int, mergeRequestID int) (comments []*models.Comment, err error) {
	comments = []*models.Comment{}

	opt := &github.PullRequestListCommentsOptions{
		ListOptions: github.ListOptions{
			PerPage: perPageItemCount,
		},
	}

	repo, _, err := s.client.Repositories.GetByID(s.ctx, int64(projectID))

	if err != nil {
		return
	}

	for {
		data, rsp, _ := s.client.PullRequests.ListComments(s.ctx, repo.Owner.GetLogin(), repo.GetName(), mergeRequestID, opt)

		comments = append(comments, s.convertGithubPullRequestsCommentsToComments(data)...)

		if rsp.NextPage == 0 {
			break
		}

		opt.Page = rsp.NextPage
	}

	return comments, nil
}

func (s *SCM) ListUsers() {

}

func (s *SCM) convertGithubPullRequestsCommentsToComments(comments []*github.PullRequestComment) []*models.Comment {
	var commentsList []*models.Comment

	for _, comment := range comments {

		//TODO: need to add more fields
		commentsList = append(commentsList, &models.Comment{
			ID:        int(comment.GetID()),
			Body:      comment.GetBody(),
			Author:    *s.convertGithubUserToUser(comment.User),
			CreatedAt: comment.GetCreatedAt(),
			UpdatedAt: comment.GetUpdatedAt(),
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
		ID:       int(user.GetID()),
		Name:     user.GetName(),
		Username: user.GetLogin(),
		State:    "",
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
		ProjectID:    int(pr.GetHead().GetRepo().GetID()),
		Title:        pr.GetTitle(),
		State:        state,
		CreatedAt:    pr.GetClosedAt(),
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

func (s *SCM) setToken(token string) error {
	s.token = token
	return nil
}

func (s *SCM) setBaseURL(baseuri string) error {
	s.baseURL = baseuri
	return nil
}
