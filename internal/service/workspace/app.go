package workspace

import (
	delivery "go-server/internal/service/workspace/delivery"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterWorkspaceService(r *gin.Engine, db *gorm.DB, httpClient *http.Client) {
	controller := delivery.NewWorkspaceController(db, httpClient)
	router := delivery.NewWorkspaceRouter(r)

	router.CreateWorkspace(controller.CreateWorkspace)
	router.DeleteWorkspace(controller.DeleteWorkspace)
	router.InviteMembers(controller.InviteMembers)
	router.RemoveMembers(controller.RemoveMembers)
	router.ChangeAccess(controller.ChangeAccess)
}
