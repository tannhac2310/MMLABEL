package banner

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

func (b *bannerService) FindBanners(ctx context.Context, opts *FindBannersOpts, limit, offset int64) ([]*repository.BannerData, *repository.CountResult, error) {
	filter := &repository.SearchBannersOpts{
		IDs:    opts.IDs,
		Name:   opts.Name,
		Limit:  limit,
		Offset: offset,
	}
	result, err := b.bannerRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	total, err := b.bannerRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	return result, total, nil
}

type FindBannersOpts struct {
	IDs  []string
	Name string
}
