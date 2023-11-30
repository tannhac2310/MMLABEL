package repository

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type InkInventoryDetailRepo interface {
	Insert(ctx context.Context, e *model.InkInventoryDetail) error
	Update(ctx context.Context, e *model.InkInventoryDetail) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchInkInventoryDetailOpts) ([]*InkInventoryDetailData, error)
	Count(ctx context.Context, s *SearchInkInventoryDetailOpts) (*CountResult, error)
}

type sInkInventoryDetailRepo struct {
}

func NewInkInventoryDetailRepo() InkInventoryDetailRepo {
	return &sInkInventoryDetailRepo{}
}

func (r *sInkInventoryDetailRepo) Insert(ctx context.Context, e *model.InkInventoryDetail) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}
