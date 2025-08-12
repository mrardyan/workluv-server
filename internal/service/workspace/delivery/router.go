package delivery

import (
	"github.com/gin-gonic/gin"
)

const (
	servicePath         = "/workspaces"
	createWorkspacePath = "/"
	deleteWorkspacePath = "/:id"
	inviteMembersPath   = "/:id/members"
	removeMembersPath   = "/:id/members/remove"
	changeAccessPath    = "/:id/members/:member_id/access"
)

type Router struct {
	routerGroup *gin.RouterGroup
}

func NewWorkspaceRouter(router *gin.Engine) *Router {
	routerGroup := router.Group(servicePath)
	return &Router{routerGroup: routerGroup}
}

func (r *Router) CreateWorkspace(handler gin.HandlerFunc) {
	r.routerGroup.POST(createWorkspacePath, handler)
}

func (r *Router) DeleteWorkspace(handler gin.HandlerFunc) {
	r.routerGroup.DELETE(deleteWorkspacePath, handler)
}

func (r *Router) InviteMembers(handler gin.HandlerFunc) {
	r.routerGroup.POST(inviteMembersPath, handler)
}

func (r *Router) RemoveMembers(handler gin.HandlerFunc) {
	r.routerGroup.POST(removeMembersPath, handler)
}

func (r *Router) ChangeAccess(handler gin.HandlerFunc) {
	r.routerGroup.PUT(changeAccessPath, handler)
}
