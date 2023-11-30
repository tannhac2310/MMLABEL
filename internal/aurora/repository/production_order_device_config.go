package repository

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type ProductionOrderDeviceConfigRepo interface {
	Insert(ctx context.Context, e *model.ProductionOrderDeviceConfig) error
	Update(ctx context.Context, e *model.ProductionOrderDeviceConfig) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchProductionOrderDeviceConfigOpts) ([]*ProductionOrderDeviceConfigData, error)
	Count(ctx context.Context, s *SearchProductionOrderDeviceConfigOpts) (*CountResult, error)
}

type sProductionOrderDeviceConfigRepo struct {
}

func NewProductionOrderDeviceConfigRepo() ProductionOrderDeviceConfigRepo {
	return &sProductionOrderDeviceConfigRepo{}
}

func (r *sProductionOrderDeviceConfigRepo) Insert(ctx context.Context, e *model.ProductionOrderDeviceConfig) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}
