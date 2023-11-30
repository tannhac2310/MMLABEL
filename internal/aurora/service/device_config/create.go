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
		DeviceID:          opt.DeviceID,
		DeviceConfig:      opt.DeviceConfig,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

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
}
