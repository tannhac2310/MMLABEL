package banner

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (b bannerService) CreateBanner(ctx context.Context, opt *CreateBannerOpts) (string, error) {
	now := time.Now()

	banner := &model.Banner{
		ID:           idutil.ULIDNow(),
		Link:         opt.Link,
		DisplayOrder: opt.DisplayOrder,
		Name:         opt.Name,
		PhotoURL:     opt.PhotoURL,
		CreatedBy:    opt.CreatedBy,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	err := b.bannerRepo.Insert(ctx, banner)
	if err != nil {
		return "", fmt.Errorf("p.bannerRepo.Insert: %w", err)
	}

	return banner.ID, nil
}

type CreateBannerOpts struct {
	Link         string
	DisplayOrder int8
	Name         string
	PhotoURL     string
	CreatedBy    string
}
