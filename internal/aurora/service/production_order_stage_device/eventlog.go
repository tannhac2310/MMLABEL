package production_order_stage_device

import (
	"context"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

type FindEventLogOpts struct {
	DeviceID string
	Date     string
}

func (p productionOrderStageDeviceService) FindEventLog(ctx context.Context, opt *FindEventLogOpts) ([]*repository.EventLogData, error) {
	return p.productionOrderStageDeviceRepo.FindEventLog(ctx, &repository.SearchEventLogOpts{
		DeviceID: opt.DeviceID,
		Date:     opt.Date,
	})
}
