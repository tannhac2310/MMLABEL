package banner

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

type Service interface {
	CreateBanner(ctx context.Context, opt *CreateBannerOpts) (string, error)
	EditBanner(ctx context.Context, opt *EditBannerOpts) error
	FindBanners(ctx context.Context, opts *FindBannersOpts, limit, offset int64) ([]*repository.BannerData, *repository.CountResult, error)
	SoftDelete(ctx context.Context, id string) error
	FindBannerByID(ctx context.Context, id string) (*repository.BannerData, error)
}

type bannerService struct {
	bannerRepo repository.BannerRepo
}

func NewService(
	bannerRepo repository.BannerRepo,
) Service {
	return &bannerService{
		bannerRepo: bannerRepo,
	}
}
