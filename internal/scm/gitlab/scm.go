package gitlab

type GitlabOption func(SCM) error

func GitlabToken(token string) GitlabOption {
	return func(s SCM) error {
		return s.setToken(token)
	}
}

func GitlabBaseURL(baseuri string) GitlabOption {
	return func(s SCM) error {
		return s.setBaseURL(baseuri)
	}
}

type SCM struct {
	token   string
	baseURL string
}

func NewSCM(options ...GitlabOption) (scm SCM, err error) {
	scm = SCM{}
	for _, fn := range options {
		if fn == nil {
			continue
		}
		if err := fn(scm); err != nil {
			return scm, err
		}
	}

	return scm, nil
}

func (s SCM) ListRepositories() {

}
func (s SCM) ListUsers() {

}

func (s SCM) setToken(token string) error {
	s.token = token
	return nil
}

func (s SCM) setBaseURL(baseuri string) error {
	s.baseURL = baseuri
	return nil
}
