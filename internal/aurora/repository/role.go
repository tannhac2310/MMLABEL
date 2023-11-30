package repository

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type RoleRepo interface {
	Insert(ctx context.Context, e *model.Role) error
	Update(ctx context.Context, e *model.Role) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchRoleOpts) ([]*RoleData, error)
	Count(ctx context.Context, s *SearchRoleOpts) (*CountResult, error)
}

type sRoleRepo struct {
}

func NewRoleRepo() RoleRepo {
	return &sRoleRepo{}
}

func (r *sRoleRepo) Insert(ctx context.Context, e *model.Role) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}
