package region

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

type Service interface {
	FindRegions(ctx context.Context, opts *FindRegionsOpts, limit, offset int64) ([]*model.Region, error)
}

type regionService struct {
	regionRepo repository.RegionRepo
}

func NewService(
	regionRepo repository.RegionRepo,
) Service {
	return &regionService{
		regionRepo: regionRepo,
	}
}
