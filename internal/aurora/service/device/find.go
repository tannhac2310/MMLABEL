package device

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

func (c *deviceService) FindDevices(ctx context.Context, opts *FindDevicesOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error) {
	filter := &repository.SearchDevicesOpts{
		Name:   opts.Name,
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	}
	devices, err := c.deviceRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	total, err := c.deviceRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*Data, 0, len(devices))
	for _, device := range devices {
		if err != nil {
			return nil, nil, err
		}
		results = append(results, &Data{
			DeviceData: device,
		})
	}
	return results, total, nil
}

type FindDevicesOpts struct {
	Name string
}
