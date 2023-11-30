package device_config

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

func (c *deviceConfigService) CreateDeviceConfig(ctx context.Context, opt *CreateDeviceConfigOpts) (string, error) {
	now := time.Now()

	deviceConfig := &model.ProductionOrderDeviceConfig{
		ID:                idutil.ULIDNow(),
		ProductionOrderID: opt.ProductionOrderID,
		DeviceID:          cockroach.String(opt.DeviceID),
		Color:             cockroach.String(opt.Color),
		Description:       cockroach.String(opt.Description),
		Search:            cockroach.String(opt.Search),
		DeviceConfig:      opt.DeviceConfig,
		CreatedBy:         opt.CreatedBy,
		CreatedAt:         now,
		UpdatedBy:         opt.CreatedBy,
		UpdatedAt:         now,
	}
	fmt.Printf("deviceConfig: %+v\n", deviceConfig)

	errTx := cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		err := c.deviceConfigRepo.Insert(ctx2, deviceConfig)
		if err != nil {
			return fmt.Errorf("c.deviceConfigRepo.Insert: %w", err)
		}

		return nil
	})
	if errTx != nil {
		return "", errTx
	}
	return deviceConfig.ID, nil
}

type CreateDeviceConfigOpts struct {
	ProductionOrderID string
	DeviceID          string
	DeviceConfig      map[string]interface{}
	Color             string
	Description       string
	Search            string
	CreatedBy         string
}
