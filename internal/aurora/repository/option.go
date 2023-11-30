package repository

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type OptionRepo interface {
	Insert(ctx context.Context, e *model.Option) error
	Update(ctx context.Context, e *model.Option) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchOptionOpts) ([]*OptionData, error)
	Count(ctx context.Context, s *SearchOptionOpts) (*CountResult, error)
}

type sOptionRepo struct {
}

func NewOptionRepo() OptionRepo {
	return &sOptionRepo{}
}

func (r *sOptionRepo) Insert(ctx context.Context, e *model.Option) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}
