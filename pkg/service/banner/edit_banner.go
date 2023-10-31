package banner

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (b *bannerService) EditBanner(ctx context.Context, opt *EditBannerOpts) error {
	var err error
	table := model.Banner{}
	updater := cockroach.NewUpdater(table.TableName(), model.BannerFieldID, opt.ID)

	updater.Set(model.BannerFieldLink, opt.Link)
	updater.Set(model.BannerFieldDisplayOrder, opt.DisplayOrder)
	updater.Set(model.BannerFieldName, opt.Name)
	updater.Set(model.BannerFieldPhotoURL, opt.PhotoURL)

	err = cockroach.UpdateFields(ctx, updater)
	if err != nil {
		return fmt.Errorf("err update banner status: %w", err)
	}

	return nil
}

type EditBannerOpts struct {
	ID           string
	Link         string
	DisplayOrder int8
	Name         string
	PhotoURL     string
}
