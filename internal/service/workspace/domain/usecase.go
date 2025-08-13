package domain

import "context"

type UseCase struct {
	Repo Repository
}

func NewUseCase(repo Repository) *UseCase {
	return &UseCase{Repo: repo}
}

func (uc *UseCase) CreateWorkspace(ctx context.Context, workspace Workspace) (Workspace, error) {
	return uc.Repo.Create(ctx, workspace)
}

func (uc *UseCase) DeleteWorkspace(ctx context.Context, workspaceID uint) error {
	return uc.Repo.Delete(ctx, workspaceID)
}

func (uc *UseCase) InviteMembers(ctx context.Context, workspaceID uint, members []Member) error {
	return uc.Repo.InviteMembers(ctx, workspaceID, members)
}

func (uc *UseCase) RemoveMembers(ctx context.Context, workspaceID uint, members []Member) error {
	return uc.Repo.RemoveMembers(ctx, workspaceID, members)
}

func (uc *UseCase) ChangeAccess(ctx context.Context, workspaceID uint, memberID uint, access Access) error {
	return uc.Repo.ChangeAccess(ctx, workspaceID, memberID, access)
}
