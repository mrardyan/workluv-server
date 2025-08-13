package delivery

import (
	"github.com/gin-gonic/gin"
)

const (
	servicePath            = "/accounts"
	createAccountPath      = "/"
	verifyEmailPath        = "/verify-email"
	resendVerificationPath = "/resend-verification"
	deleteAccountPath      = "/:id"
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

func (r *Router) VerifyEmail(handler gin.HandlerFunc) {
	r.routerGroup.GET(verifyEmailPath, handler)
	r.routerGroup.POST(verifyEmailPath, handler)
}

func (r *Router) ResendVerificationEmail(handler gin.HandlerFunc) {
	r.routerGroup.POST(resendVerificationPath, handler)
}

func (r *Router) DeleteAccount(handler gin.HandlerFunc) {
	r.routerGroup.DELETE(deleteAccountPath, handler)
}
