package delivery

import (
	"github.com/gin-gonic/gin"
)

const (
	servicePath    = "/users"
	createUserPath = "/"
	deleteUserPath = "/:id"
)

type Router struct {
	routerGroup *gin.RouterGroup
}

func NewUserRouter(router *gin.Engine) *Router {
	routerGroup := router.Group(servicePath)
	return &Router{routerGroup: routerGroup}
}

func (r *Router) CreateUser(handler gin.HandlerFunc) {
	r.routerGroup.POST(createUserPath, handler)
}

func (r *Router) DeleteUser(handler gin.HandlerFunc) {
	r.routerGroup.DELETE(deleteUserPath, handler)
}
