package selfmerging

type SCM interface {
	ListRepositories()
	ListUsers()
}

type SelfMergingService interface {
	Calculate()
}

type selfMerging struct {
	scm SCM
}

func NewSelfMergingService(scm SCM) SelfMergingService {
	return selfMerging{scm: scm}
}

func (tor selfMerging) Calculate() {

}
