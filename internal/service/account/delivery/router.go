package delivery

import (
	"github.com/gin-gonic/gin"
)

const (
	servicePath       = "/accounts"
	createAccountPath = "/"
	deleteAccountPath = "/:id"
)

type Router struct {
	routerGroup *gin.RouterGroup
}

func NewAccountRouter(router *gin.Engine) *Router {
	routerGroup := router.Group(servicePath)
	return &Router{routerGroup: routerGroup}
}

func (r *Router) CreateAccount(handler gin.HandlerFunc) {
	r.routerGroup.POST(createAccountPath, handler)
}

func (r *Router) DeleteAccount(handler gin.HandlerFunc) {
	r.routerGroup.DELETE(deleteAccountPath, handler)
}
