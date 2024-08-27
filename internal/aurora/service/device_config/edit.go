package device_config

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

func (c *deviceConfigService) EditDeviceConfig(ctx context.Context, opt *EditDeviceConfigOpts) error {
	var err error
	table := model.ProductionOrderDeviceConfig{}
	updater := cockroach.NewUpdater(table.TableName(), model.ProductionOrderDeviceConfigFieldID, opt.ID)

	updater.Set(model.ProductionOrderDeviceConfigFieldProductionOrderID, opt.ProductionOrderID)
	updater.Set(model.ProductionOrderDeviceConfigFieldProductionPlanID, opt.ProductionPlanID)
	updater.Set(model.ProductionOrderDeviceConfigFieldDeviceID, opt.DeviceID)
	updater.Set(model.ProductionOrderDeviceConfigFieldDeviceConfig, opt.DeviceConfig)
	updater.Set(model.ProductionOrderDeviceConfigFieldColor, opt.Color)
	updater.Set(model.ProductionOrderDeviceConfigFieldDescription, opt.Description)
	updater.Set(model.ProductionOrderDeviceConfigFieldSearch, opt.Search)

	updater.Set(model.ProductionOrderDeviceConfigFieldUpdatedAt, time.Now())
	updater.Set(model.ProductionOrderDeviceConfigFieldUpdatedBy, opt.UpdatedBy)

	err = cockroach.UpdateFields(ctx, updater)
	if err != nil {
		return fmt.Errorf("update deviceConfig failed %w", err)
	}
	return nil
}

type EditDeviceConfigOpts struct {
	ID                string
	ProductionOrderID string
	ProductionPlanID  string
	DeviceID          string
	DeviceConfig      map[string]interface{}
	Color             string
	Description       string
	Search            string
	UpdatedBy         string
}
