package delivery

import (
	"database/sql"
	"net/http"
	"strconv"
	"workluv/internal/service/workspace/domain"
	"workluv/internal/service/workspace/repository"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	UseCase    *domain.UseCase
	HTTPClient *http.Client
}

func NewWorkspaceController(db *sql.DB, httpClient *http.Client) *Controller {
	repo := repository.NewRepository(db)
	useCase := domain.NewUseCase(repo)
	return &Controller{
		UseCase:    useCase,
		HTTPClient: httpClient,
	}
}

func (ctl *Controller) CreateWorkspace(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		OwnerID   uint   `json:"owner_id" binding:"required"`
		OwnerName string `json:"owner_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	workspace, err := ctl.UseCase.CreateWorkspace(c.Request.Context(), domain.Workspace{
		Name: req.Name,
		Owner: domain.Owner{
			ID:       req.OwnerID,
			Username: req.OwnerName,
		},
		Members: []domain.Member{},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, workspace)
}

func (ctl *Controller) DeleteWorkspace(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	err = ctl.UseCase.DeleteWorkspace(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workspace deleted successfully"})
}

func (ctl *Controller) InviteMembers(c *gin.Context) {
	workspaceID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	var req struct {
		Members []struct {
			ID       uint   `json:"id" binding:"required"`
			Username string `json:"username" binding:"required"`
			Access   string `json:"access" binding:"required"`
		} `json:"members" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	members := make([]domain.Member, len(req.Members))
	for i, m := range req.Members {
		members[i] = domain.Member{
			ID:       m.ID,
			Username: m.Username,
			Access:   domain.Access(m.Access),
		}
	}

	err = ctl.UseCase.InviteMembers(c.Request.Context(), uint(workspaceID), members)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Members invited successfully"})
}

func (ctl *Controller) RemoveMembers(c *gin.Context) {
	workspaceID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	var req struct {
		Members []struct {
			ID       uint   `json:"id" binding:"required"`
			Username string `json:"username" binding:"required"`
		} `json:"members" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	members := make([]domain.Member, len(req.Members))
	for i, m := range req.Members {
		members[i] = domain.Member{
			ID:       m.ID,
			Username: m.Username,
		}
	}

	err = ctl.UseCase.RemoveMembers(c.Request.Context(), uint(workspaceID), members)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Members removed successfully"})
}

func (ctl *Controller) ChangeAccess(c *gin.Context) {
	workspaceID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	memberID, err := strconv.ParseUint(c.Param("member_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid member ID"})
		return
	}

	var req struct {
		Access string `json:"access" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err = ctl.UseCase.ChangeAccess(c.Request.Context(), uint(workspaceID), uint(memberID), domain.Access(req.Access))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Access changed successfully"})
}
