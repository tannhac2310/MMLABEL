package device

import (
	"context"
	"fmt"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

func (c *deviceService) EditDevice(ctx context.Context, opt *EditDeviceOpts) error {
	var err error
	table := model.Device{}
	updater := cockroach.NewUpdater(table.TableName(), model.DeviceFieldID, opt.ID)

	updater.Set(model.DeviceFieldName, opt.Name)
	updater.Set(model.DeviceFieldName, opt.Name)
	updater.Set(model.DeviceFieldCode, opt.Code)
	updater.Set(model.DeviceFieldOptionID, opt.OptionID)
	updater.Set(model.DeviceFieldStatus, opt.Status)
	updater.Set(model.DeviceFieldData, opt.Data)

	updater.Set(model.DeviceFieldUpdatedAt, time.Now())

	err = cockroach.UpdateFields(ctx, updater)
	if err != nil {
		return fmt.Errorf("update device failed %w", err)
	}
	return nil
}

type EditDeviceOpts struct {
	ID       string
	Name     string
	Code     string
	OptionID string
	Status   enum.CommonStatus
	Data     map[string]interface{}
}
