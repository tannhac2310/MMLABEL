package production_order_stage_device

import (
	"context"
	"fmt"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"
)

func (p productionOrderStageDeviceService) CalcOEE(ctx context.Context, dateFrom string, dateTo string) (map[string]model.OEE, error) {
	now := time.Now()
	stringNow := now.Format("2006-01-02")

	if dateFrom == "" {
		dateFrom = stringNow
	}
	if dateTo == "" {
		dateTo = stringNow
	}

	_, err := time.Parse("2006-01-02", dateFrom)
	if err != nil {
		return nil, fmt.Errorf("invalid dateFrom format: expected YYYY-MM-DD, got %s", dateFrom)
	}

	_, err = time.Parse("2006-01-02", dateTo)
	if err != nil {
		return nil, fmt.Errorf("invalid dateTo format: expected YYYY-MM-DD, got %s", dateTo)
	}

	listDeviceProgressStatusHistory, err := p.sDeviceProgressStatusHistoryRepo.FindByDate(ctx, dateFrom, dateTo)
	if err != nil {
		return nil, fmt.Errorf("p.CalcOEE: %w", err)
	}

	assignedTime, err := p.productionOrderStageDeviceRepo.GetAssignedTimeByDate(ctx, dateFrom, dateTo)
	if err != nil {
		return nil, fmt.Errorf("p.CalcOEE: %w", err)
	}

	result := make(map[string]model.OEE)
	var lastDeviceProgressStatusHistory *repository.DeviceProgressStatusHistoryData

	for _, deviceProgressStatusHistory := range listDeviceProgressStatusHistory {
		tmp, exists := result[deviceProgressStatusHistory.DeviceID]
		if !exists {
			tmp = model.OEE{
				ActualWorkingTime:  28800,
				JobRunningTime:     0,
				AssignedWorkTime:   assignedTime[deviceProgressStatusHistory.DeviceID].TotalRuntime,
				TotalQuantity:      assignedTime[deviceProgressStatusHistory.DeviceID].TotalQuantity,
				TotalDefective:     assignedTime[deviceProgressStatusHistory.DeviceID].TotalDefective,
				DowntimeStatistics: make(map[string]string),
			}
			lastDeviceProgressStatusHistory = &deviceProgressStatusHistory
			result[deviceProgressStatusHistory.DeviceID] = tmp
			continue
		}
		if lastDeviceProgressStatusHistory == nil {
			continue
		}
		duration := deviceProgressStatusHistory.CreatedAt.Sub(lastDeviceProgressStatusHistory.CreatedAt).Seconds()

		if deviceProgressStatusHistory.ProcessStatus == enum.ProductionOrderStageDeviceStatusStart &&
			lastDeviceProgressStatusHistory.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed {
			tmp.Downtime += int64(duration)
		}

		if deviceProgressStatusHistory.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed ||
			deviceProgressStatusHistory.ProcessStatus == enum.ProductionOrderStageDeviceStatusComplete {
			tmp.JobRunningTime += int64(duration)

			if deviceProgressStatusHistory.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed {
				tmp.DowntimeStatistics[deviceProgressStatusHistory.CreatedAt.Format(time.RFC3339)] = deviceProgressStatusHistory.ErrorReason.String
			}
		}

		lastDeviceProgressStatusHistory = &deviceProgressStatusHistory
		result[deviceProgressStatusHistory.DeviceID] = tmp
	}
	return result, nil
}
