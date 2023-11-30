package repository

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type UsernamePasswordRepo interface {
	Insert(ctx context.Context, e *model.UsernamePassword) error
	Update(ctx context.Context, e *model.UsernamePassword) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchUsernamePasswordOpts) ([]*UsernamePasswordData, error)
	Count(ctx context.Context, s *SearchUsernamePasswordOpts) (*CountResult, error)
}

type sUsernamePasswordRepo struct {
}

func NewUsernamePasswordRepo() UsernamePasswordRepo {
	return &sUsernamePasswordRepo{}
}

func (r *sUsernamePasswordRepo) Insert(ctx context.Context, e *model.UsernamePassword) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}
