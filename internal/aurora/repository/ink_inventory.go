package repository

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type InkInventoryRepo interface {
	Insert(ctx context.Context, e *model.InkInventory) error
	Update(ctx context.Context, e *model.InkInventory) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchInkInventoryOpts) ([]*InkInventoryData, error)
	Count(ctx context.Context, s *SearchInkInventoryOpts) (*CountResult, error)
}

type sInkInventoryRepo struct {
}

func NewInkInventoryRepo() InkInventoryRepo {
	return &sInkInventoryRepo{}
}

func (r *sInkInventoryRepo) Insert(ctx context.Context, e *model.InkInventory) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}
