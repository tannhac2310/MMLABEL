package production_order_stage_device

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

// parseDateOrDefault - Helper function to parse date or use default
func parseDateOrDefault(date, defaultDate string) (string, error) {
	if date == "" {
		return defaultDate, nil
	}
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return "", fmt.Errorf("invalid date format: expected YYYY-MM-DD, got %s", date)
	}
	return date, nil
}

func (p productionOrderStageDeviceService) CalcOEE(ctx context.Context, dateFrom, dateTo string) (map[string]model.OEE, error) {
	now := time.Now().Format("2006-01-02")

	// Parse dates with default values
	var err error
	if dateFrom, err = parseDateOrDefault(dateFrom, now); err != nil {
		return nil, err
	}
	if dateTo, err = parseDateOrDefault(dateTo, now); err != nil {
		return nil, err
	}

	// Fetch necessary data
	listDeviceProgressStatusHistory, err := p.sDeviceProgressStatusHistoryRepo.FindByDate(ctx, dateFrom, dateTo)
	if err != nil {
		return nil, fmt.Errorf("p.CalcOEE: %w", err)
	}

	assignedWork, err := p.productionOrderStageDeviceRepo.GetAssignedByDate(ctx, dateFrom, dateTo)
	if err != nil {
		return nil, fmt.Errorf("p.CalcOEE: %w", err)
	}
	fmt.Println(assignedWork)
	// Initialize result map
	result := make(map[string]model.OEE)

	var lastHistory *repository.DeviceProgressStatusHistoryData

	for i := range listDeviceProgressStatusHistory {
		history := &listDeviceProgressStatusHistory[i]
		deviceID := history.DeviceID

		oee, exists := result[deviceID]
		if !exists {
			oee = model.OEE{
				ActualWorkingTime:  28800,
				DowntimeStatistics: make(map[string]string),
				AssignedWork:       assignedWork[deviceID],
			}

			for _, assigned := range assignedWork[deviceID] {
				oee.TotalQuantity += assigned.Quantity
				if assigned.Settings != nil {
					if val, ok := assigned.Settings["san_pham_loi"].(int64); ok {
						oee.TotalDefective += val
					}
				}
				oee.AssignedWorkTime += assigned.EstimatedCompleteAt.Time.Unix() - assigned.EstimatedStartAt.Time.Unix()
			}

		}

		if lastHistory != nil {
			duration := history.CreatedAt.Sub(lastHistory.CreatedAt).Seconds()

			if history.ProcessStatus == enum.ProductionOrderStageDeviceStatusStart &&
				lastHistory.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed {
				oee.Downtime += int64(duration)
			}

			if history.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed ||
				history.ProcessStatus == enum.ProductionOrderStageDeviceStatusComplete {
				oee.JobRunningTime += int64(duration)

				if history.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed {
					oee.DowntimeStatistics[history.CreatedAt.Format(time.RFC3339)] = history.ErrorReason.String
				}
			}

		}

		result[deviceID] = oee
		lastHistory = history
	}

	return result, nil
}
