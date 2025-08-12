package domain

type UseCase struct {
	Repo Repository
}

func NewUseCase(repo Repository) *UseCase {
	return &UseCase{Repo: repo}
}

func (uc *UseCase) CreateUser(user User) (User, error) {
	return uc.Repo.Create(user)
}

func (uc *UseCase) DeleteUser(id uint) error {
	return uc.Repo.Delete(id)
}
