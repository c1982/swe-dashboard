package github

type GithubOption func(SCM) error

type SCM struct {
}

func NewSCM(options ...GithubOption) SCM {
	return SCM{}
}

func (s SCM) ListRepositories() {

}

func (s SCM) ListUsers() {

}
