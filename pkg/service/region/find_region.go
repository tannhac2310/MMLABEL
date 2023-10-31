package region

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

func (b *regionService) FindRegions(ctx context.Context, opts *FindRegionsOpts, limit, offset int64) ([]*model.Region, error) {
	return b.regionRepo.Search(ctx, &repository.SearchRegionsOpts{
		IDs:      opts.IDs,
		Name:     opts.Name,
		Search:   opts.Search,
		ParentID: opts.ParentID,
		Limit:    limit,
		Offset:   offset,
	})
}

type FindRegionsOpts struct {
	IDs      []int64
	Name     string
	Search   string
	ParentID int64
}
