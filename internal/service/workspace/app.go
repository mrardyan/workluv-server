package workspace

import (
	"database/sql"
	"net/http"
	delivery "workluv/internal/service/workspace/delivery"

	"github.com/gin-gonic/gin"
)

func RegisterWorkspaceService(r *gin.Engine, db *sql.DB, httpClient *http.Client) {
	controller := delivery.NewWorkspaceController(db, httpClient)
	router := delivery.NewWorkspaceRouter(r)

	router.CreateWorkspace(controller.CreateWorkspace)
	router.DeleteWorkspace(controller.DeleteWorkspace)
	router.InviteMembers(controller.InviteMembers)
	router.RemoveMembers(controller.RemoveMembers)
	router.ChangeAccess(controller.ChangeAccess)
}
