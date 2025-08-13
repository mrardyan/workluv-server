package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go-server/internal/service/workspace/domain"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) domain.Repository {
	return &Repository{DB: db}
}

func (r *Repository) Create(ctx context.Context, workspace domain.Workspace) (domain.Workspace, error) {
	query := `
		INSERT INTO workspaces (name)
		VALUES ($1)
		RETURNING id, name`

	row := r.DB.QueryRowContext(ctx, query, workspace.Name)

	var createdWorkspace domain.Workspace
	err := row.Scan(
		&createdWorkspace.ID,
		&createdWorkspace.Name,
	)
	if err != nil {
		return domain.Workspace{}, fmt.Errorf("failed to create workspace: %w", err)
	}

	// Copy other fields that aren't stored in database
	createdWorkspace.Owner = workspace.Owner
	createdWorkspace.Members = workspace.Members

	return createdWorkspace, nil
}

func (r *Repository) Delete(ctx context.Context, workspaceID uint) error {
	query := `DELETE FROM workspaces WHERE id = $1`
	result, err := r.DB.ExecContext(ctx, query, workspaceID)
	if err != nil {
		return fmt.Errorf("failed to delete workspace: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("workspace not found")
	}

	return nil
}

func (r *Repository) InviteMembers(ctx context.Context, workspaceID uint, members []domain.Member) error {
	// TODO: Implement member invitation logic using SQL
	// This would typically involve creating member records and sending invitations
	return nil
}

func (r *Repository) RemoveMembers(ctx context.Context, workspaceID uint, members []domain.Member) error {
	// TODO: Implement member removal logic using SQL
	// This would typically involve removing member records and updating access
	return nil
}

func (r *Repository) ChangeAccess(ctx context.Context, workspaceID uint, memberID uint, access domain.Access) error {
	// TODO: Implement access change logic using SQL
	// This would typically involve updating member access permissions
	return nil
}
