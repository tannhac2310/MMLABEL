package banner

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

func (b *bannerService) FindBannerByID(ctx context.Context, id string) (*repository.BannerData, error) {
	banners, _, err := b.FindBanners(ctx, &FindBannersOpts{
		IDs: []string{id},
	}, 1, 0)

	if err != nil {
		return nil, err
	}
	if len(banners) != 1 {
		return nil, fmt.Errorf("banner.Search:FindBannerByID not found")
	}

	return banners[0], nil
}
