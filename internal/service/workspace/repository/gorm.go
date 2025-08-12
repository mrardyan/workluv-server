package repository

import (
	"go-server/internal/service/workspace/domain"

	"gorm.io/gorm"
)

type GormRepository struct {
	DB *gorm.DB
}

func NewGormRepository(db *gorm.DB) domain.Repository {
	return &GormRepository{DB: db}
}

func (r *GormRepository) Create(workspace domain.Workspace) (domain.Workspace, error) {
	err := r.DB.Create(&workspace).Error
	return workspace, err
}

func (r *GormRepository) Delete(workspaceID uint) error {
	return r.DB.Delete(&domain.Workspace{}, workspaceID).Error
}

func (r *GormRepository) InviteMembers(workspaceID uint, members []domain.Member) error {
	// TODO: Implement member invitation logic
	// This would typically involve creating member records and sending invitations
	return nil
}

func (r *GormRepository) RemoveMembers(workspaceID uint, members []domain.Member) error {
	// TODO: Implement member removal logic
	// This would typically involve removing member records and updating access
	return nil
}

func (r *GormRepository) ChangeAccess(workspaceID uint, memberID uint, access domain.Access) error {
	// TODO: Implement access change logic
	// This would typically involve updating member access permissions
	return nil
}
