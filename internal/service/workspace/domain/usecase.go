package domain

type UseCase struct {
	Repo Repository
}

func NewUseCase(repo Repository) *UseCase {
	return &UseCase{Repo: repo}
}

func (uc *UseCase) CreateWorkspace(workspace Workspace) (Workspace, error) {
	return uc.Repo.Create(workspace)
}

func (uc *UseCase) DeleteWorkspace(workspaceID uint) error {
	return uc.Repo.Delete(workspaceID)
}

func (uc *UseCase) InviteMembers(workspaceID uint, members []Member) error {
	return uc.Repo.InviteMembers(workspaceID, members)
}

func (uc *UseCase) RemoveMembers(workspaceID uint, members []Member) error {
	return uc.Repo.RemoveMembers(workspaceID, members)
}

func (uc *UseCase) ChangeAccess(workspaceID uint, memberID uint, access Access) error {
	return uc.Repo.ChangeAccess(workspaceID, memberID, access)
}
