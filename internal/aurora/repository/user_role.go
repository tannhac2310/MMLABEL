package repository

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type UserRoleRepo interface {
	Insert(ctx context.Context, e *model.UserRole) error
	Update(ctx context.Context, e *model.UserRole) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchUserRoleOpts) ([]*UserRoleData, error)
	Count(ctx context.Context, s *SearchUserRoleOpts) (*CountResult, error)
}

type sUserRoleRepo struct {
}

func NewUserRoleRepo() UserRoleRepo {
	return &sUserRoleRepo{}
}

func (r *sUserRoleRepo) Insert(ctx context.Context, e *model.UserRole) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}
