package repository

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type UserRepo interface {
	Insert(ctx context.Context, e *model.User) error
	Update(ctx context.Context, e *model.User) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchUserOpts) ([]*UserData, error)
	Count(ctx context.Context, s *SearchUserOpts) (*CountResult, error)
}

type sUserRepo struct {
}

func NewUserRepo() UserRepo {
	return &sUserRepo{}
}

func (r *sUserRepo) Insert(ctx context.Context, e *model.User) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}
