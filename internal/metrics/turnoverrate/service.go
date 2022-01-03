package turnoverrate

type SCM interface {
	ListRepositories()
	ListUsers()
}

type TurnOverrateService interface {
	Calculate()
}

type turnOverRate struct {
	scm SCM
}

func NewTurnOverRate(scm SCM) TurnOverrateService {
	return turnOverRate{scm: scm}
}

func (tor turnOverRate) Calculate() {

}
