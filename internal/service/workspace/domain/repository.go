package domain

import "context"

type Repository interface {
	Create(ctx context.Context, workspace Workspace) (Workspace, error)
	Delete(ctx context.Context, workspaceID uint) error
	InviteMembers(ctx context.Context, workspaceID uint, members []Member) error
	RemoveMembers(ctx context.Context, workspaceID uint, members []Member) error
	ChangeAccess(ctx context.Context, workspaceID uint, memberID uint, access Access) error
}
