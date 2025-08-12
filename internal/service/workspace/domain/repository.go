package domain

type Repository interface {
	Create(workspace Workspace) (Workspace, error)
	Delete(workspaceID uint) error
	InviteMembers(workspaceID uint, members []Member) error
	RemoveMembers(workspaceID uint, members []Member) error
	ChangeAccess(workspaceID uint, memberID uint, access Access) error
}
