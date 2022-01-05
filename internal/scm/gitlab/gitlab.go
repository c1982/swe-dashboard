package gitlab

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
