package delivery

import (
	"github.com/gin-gonic/gin"
)

const (
	servicePath      = "/auth"
	loginPath        = "/login"
	logoutPath       = "/logout"
	logoutAllPath    = "/logout-all"
	refreshTokenPath = "/refresh"
)

type Router struct {
	routerGroup *gin.RouterGroup
}

// NewSessionRouter creates a new session router
func NewSessionRouter(router *gin.Engine) *Router {
	routerGroup := router.Group(servicePath)
	return &Router{routerGroup: routerGroup}
}

// Login registers the login endpoint
func (r *Router) Login(handler gin.HandlerFunc) {
	r.routerGroup.POST(loginPath, handler)
}

// Logout registers the logout endpoint
func (r *Router) Logout(handler gin.HandlerFunc) {
	r.routerGroup.POST(logoutPath, handler)
}

// LogoutAll registers the logout all endpoint
func (r *Router) LogoutAll(handler gin.HandlerFunc) {
	r.routerGroup.POST(logoutAllPath, handler)
}

// RefreshToken registers the refresh token endpoint
func (r *Router) RefreshToken(handler gin.HandlerFunc) {
	r.routerGroup.POST(refreshTokenPath, handler)
}
