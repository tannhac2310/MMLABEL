package repository

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type ProductQualityRepo interface {
	Insert(ctx context.Context, e *model.ProductQuality) error
	Update(ctx context.Context, e *model.ProductQuality) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchProductQualityOpts) ([]*ProductQualityData, error)
	Count(ctx context.Context, s *SearchProductQualityOpts) (*CountResult, error)
}

type sProductQualityRepo struct {
}

func NewProductQualityRepo() ProductQualityRepo {
	return &sProductQualityRepo{}
}

func (r *sProductQualityRepo) Insert(ctx context.Context, e *model.ProductQuality) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}
