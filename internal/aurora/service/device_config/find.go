package device_config

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
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
	inkIds := make([]string, 0, len(deviceConfigs))
	for _, deviceConfig := range deviceConfigs {
		if deviceConfig.InkID.String != "" {
			inkIds = append(inkIds, deviceConfig.InkID.String)
		}
	}

	inkData, err := c.inkRepo.Search(ctx, &repository.SearchInkOpts{
		IDs:    inkIds,
		Limit:  1000,
		Offset: 0,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("lấy danh sách mực thất bại: %w", err)
	}

	inkMap := make(map[string]*repository.InkData, len(inkData))
	for _, ink := range inkData {
		inkMap[ink.ID] = ink
	}

	results := make([]*Data, 0, len(deviceConfigs))
	for _, deviceConfig := range deviceConfigs {
		tenMauMuc := ""
		tenLoaiMuc := ""
		maMauMuc := ""
		if deviceConfig.InkID.String != "" {
			ink, ok := inkMap[deviceConfig.InkID.String]
			if ok {
				tenLoaiMuc = ink.LoaiMuc
				tenMauMuc = ink.Name
				maMauMuc = ink.Code
			}
		}
		deviceConfig.TenMauMuc = cockroach.String(tenMauMuc)
		deviceConfig.TenLoaiMuc = cockroach.String(tenLoaiMuc)
		deviceConfig.MaMauMuc = cockroach.String(maMauMuc)

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
