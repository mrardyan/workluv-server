package domain

type UseCase struct {
	Repo Repository
}

func NewUseCase(repo Repository) *UseCase {
	return &UseCase{Repo: repo}
}

func (uc *UseCase) CreateAccount(account Account) (Account, error) {
	return uc.Repo.Create(account)
}

func (uc *UseCase) DeleteAccount(id uint) error {
	return uc.Repo.Delete(id)
}
