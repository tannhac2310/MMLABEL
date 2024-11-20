package device

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

func (c *deviceService) CreateDevice(ctx context.Context, opt *CreateDeviceOpts) (string, error) {
	now := time.Now()

	device := &model.Device{
		ID:       idutil.ULIDNow(),
		Name:     opt.Name,
		Code:     opt.Code,
		OptionID: cockroach.String(opt.OptionID),
		//Data:      opt.Data,
		Status:    opt.Status,
		CreatedBy: opt.CreatedBy,
		CreatedAt: now,
		UpdatedAt: now,
	}

	errTx := cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		err := c.deviceRepo.Insert(ctx2, device)
		if err != nil {
			return fmt.Errorf("c.deviceRepo.Insert: %w", err)
		}

		return nil
	})
	if errTx != nil {
		return "", errTx
	}
	return device.ID, nil
}

type CreateDeviceOpts struct {
	Name     string
	Code     string
	OptionID string
	//Data      map[string]interface{}
	Status    enum.CommonStatus
	CreatedBy string
}
