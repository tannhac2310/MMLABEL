package production_order_stage_device

import (
	"context"
	"fmt"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"
)

type FindLostTimeOpts struct {
	DeviceID string
	Date     string // format: 2006-01-02
}

func (p productionOrderStageDeviceService) CalculateLostTime(ctx context.Context, opt *FindLostTimeOpts) (float64, error) {
	// find FindProcessDeviceHistory by deviceID and start of date and end of date
	createdFrom, _ := time.Parse("2006-01-02 15:04:05", opt.Date+" 00:00:00")
	createdTo, _ := time.Parse("2006-01-02 15:04:05", opt.Date+" 23:59:59")
	processDeviceHistory, _, err := p.FindProcessDeviceHistory(ctx, &FindProcessDeviceHistoryOpts{
		DeviceID:      opt.DeviceID,
		CreatedFrom:   createdFrom,
		CreatedTo:     createdTo,
		ProcessStatus: []int8{int8(enum.ProductionOrderStageDeviceStatusFailed)},
	}, &repository.Sort{
		Order: "ASC",
		By:    "created_at",
	}, 100000, 0)
	if err != nil {
		return 0, fmt.Errorf("p.FindProcessDeviceHistory: %w", err)
	}
	// calculate lost time
	var lostTime float64
	now := time.Now()
	for i, processDevice := range processDeviceHistory {
		if i == 0 {
			continue
		}
		if processDevice.UpdatedAt.Valid && processDevice.IsResolved == 1 {
			lostTime += processDevice.UpdatedAt.Time.Sub(processDeviceHistory[i-1].CreatedAt).Minutes()
		} else {
			lostTime += now.Sub(processDeviceHistory[i-1].CreatedAt).Minutes()
		}
	}
	lunchTimeAt, _ := time.Parse("15:04:05", "11:30:00")
	lunchTimeTo, _ := time.Parse("15:04:05", "12:00:00")
	if now.After(lunchTimeAt) && now.Before(lunchTimeTo) {
		lostTime -= now.Sub(lunchTimeAt).Minutes()
	} else {
		lostTime -= 30
	}

	if lostTime < 0 {
		lostTime = 0
	}

	return lostTime, nil
}
