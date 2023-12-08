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

type AvailabilityTime struct {
	AvailabilityTime int64 `json:"availabilityTime"`
	LossTime         int64 `json:"lossTime"`
}

func (p productionOrderStageDeviceService) FindAvailabilityTime(ctx context.Context, opt *FindLostTimeOpts) (*AvailabilityTime, error) {
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
		return nil, fmt.Errorf("p.FindProcessDeviceHistory: %w", err)
	}
	// calculate lost time
	var lossTime float64
	now := time.Now()
	for i, processDevice := range processDeviceHistory {
		if i == 0 {
			continue
		}
		if processDevice.UpdatedAt.Valid && processDevice.IsResolved == 1 {
			lossTime += processDevice.UpdatedAt.Time.Sub(processDeviceHistory[i-1].CreatedAt).Minutes()
		} else {
			lossTime += now.Sub(processDeviceHistory[i-1].CreatedAt).Minutes()
		}
	}
	lunchTimeAt, _ := time.Parse("15:04:05", "11:30:00")
	lunchTimeTo, _ := time.Parse("15:04:05", "12:00:00")
	if now.After(lunchTimeAt) && now.Before(lunchTimeTo) {
		lossTime -= now.Sub(lunchTimeAt).Minutes()
	} else {
		lossTime -= 30
	}

	if lossTime < 0 {
		lossTime = 0
	}
	var availabilityTime int64
	// find availability time from device working history
	availabilityTimes, err := p.sDeviceWorkingHistoryRepo.Search(ctx, &repository.SearchDeviceWorkingHistoryOpts{
		DeviceID: opt.DeviceID,
		Date:     opt.Date,
		Limit:    1,
		Offset:   0,
	})
	if err != nil {
		// do nothing
	}
	if len(availabilityTimes) == 1 {
		availabilityTime = availabilityTimes[0].WorkingTime
	}
	return &AvailabilityTime{
		AvailabilityTime: availabilityTime,
		LossTime:         int64(lossTime),
	}, nil
}
