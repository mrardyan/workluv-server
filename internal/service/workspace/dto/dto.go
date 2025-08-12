package dto

// CreateWorkspaceRequest represents the request for creating a workspace
type CreateWorkspaceRequest struct {
	Name      string `json:"name" binding:"required"`
	OwnerID   uint   `json:"owner_id" binding:"required"`
	OwnerName string `json:"owner_name" binding:"required"`
}

// CreateWorkspaceResponse represents the response after creating a workspace
type CreateWorkspaceResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	OwnerID uint   `json:"owner_id"`
	Message string `json:"message"`
}

// InviteMembersRequest represents the request for inviting members to a workspace
type InviteMembersRequest struct {
	Members []MemberInvite `json:"members" binding:"required"`
}

type MemberInvite struct {
	ID       uint   `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Access   string `json:"access" binding:"required"`
}

// RemoveMembersRequest represents the request for removing members from a workspace
type RemoveMembersRequest struct {
	Members []MemberRemove `json:"members" binding:"required"`
}

type MemberRemove struct {
	ID       uint   `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
}

// ChangeAccessRequest represents the request for changing member access
type ChangeAccessRequest struct {
	Access string `json:"access" binding:"required"`
}

// GenericResponse represents a generic API response
type GenericResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
