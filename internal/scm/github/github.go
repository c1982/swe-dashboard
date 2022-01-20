package github

type GithubOption func(*SCM) error

func GithubToken(token string) GithubOption {
	return func(s *SCM) error {
		return s.setToken(token)
	}
}

func GithubBaseURL(baseuri string) GithubOption {
	return func(s *SCM) error {
		return s.setBaseURL(baseuri)
	}
}

func GithubOrganizations(list []string) GithubOption {
	return func(s *SCM) error {
		return s.setOrganizations(list)
	}
}
