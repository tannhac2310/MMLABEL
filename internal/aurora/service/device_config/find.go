package device_config

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

func (c *deviceConfigService) FindDeviceConfigs(ctx context.Context, opts *FindDeviceConfigsOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error) {
	filter := &repository.SearchProductionOrderDeviceConfigOpts{
		Search: opts.Search,
		Limit:  limit,
		ProductionOrderID: opts.ProductionOrderID,
		Offset: offset,
		Sort:   sort,
	}
	deviceConfigs, err := c.deviceConfigRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	total, err := c.deviceConfigRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*Data, 0, len(deviceConfigs))
	for _, deviceConfig := range deviceConfigs {
		if err != nil {
			return nil, nil, err
		}

		results = append(results, &Data{
			ProductionOrderDeviceConfigData: deviceConfig,
		})
	}
	return results, total, nil
}

type FindDeviceConfigsOpts struct {
	Search string
	ProductionOrderID string
}
