package office

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

func (b *officeService) FindOffices(ctx context.Context, opts *FindOfficesOpts, limit, offset int64) ([]*repository.OfficeData, *repository.CountResult, error) {
	filter := &repository.SearchOfficesOpts{
		IDs:    opts.IDs,
		Name:   opts.Name,
		Search: opts.Search,
		Limit:  limit,
		Offset: offset,
	}
	offices, err := b.officeRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	total, err := b.officeRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	return offices, total, nil
}

type FindOfficesOpts struct {
	IDs        []string
	Name       string
	Search     string
	Address    string
	ProvinceID string
	DistrictID string
	Status     enum.CommonStatus
}
