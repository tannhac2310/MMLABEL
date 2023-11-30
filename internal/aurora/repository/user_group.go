package repository

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type UserGroupRepo interface {
	Insert(ctx context.Context, e *model.UserGroup) error
	Update(ctx context.Context, e *model.UserGroup) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchUserGroupOpts) ([]*UserGroupData, error)
	Count(ctx context.Context, s *SearchUserGroupOpts) (*CountResult, error)
}

type sUserGroupRepo struct {
}

func NewUserGroupRepo() UserGroupRepo {
	return &sUserGroupRepo{}
}

func (r *sUserGroupRepo) Insert(ctx context.Context, e *model.UserGroup) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}
