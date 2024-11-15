package device_config

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

func (c *deviceConfigService) FindDeviceConfigs(ctx context.Context, opts *FindDeviceConfigsOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error) {
	filter := &repository.SearchProductionOrderDeviceConfigOpts{
		IDs:               opts.IDs,
		Search:            opts.Search,
		ProductionOrderID: opts.ProductionOrderID,
		ProductionPlanID:  opts.ProductionPlanID,
		DeviceType:        opts.DeviceType,
		Limit:             limit,
		Offset:            offset,
		Sort:              sort,
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
		results = append(results, &Data{
			ProductionOrderDeviceConfigData: deviceConfig,
		})
	}
	return results, total, nil
}

type FindDeviceConfigsOpts struct {
	IDs               []string
	Search            string
	ProductionOrderID string
	ProductionPlanID  string
	DeviceType        string
}
