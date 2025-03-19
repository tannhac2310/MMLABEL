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

func (p productionOrderStageDeviceService) CalcOEEByDevice(ctx context.Context, dateFrom, dateTo string) (map[string]model.OEE, error) {
	now := time.Now().Format("2006-01-02")

	var err error
	if dateFrom, err = parseDateOrDefault(dateFrom, now); err != nil {
		return nil, err
	}
	if dateTo, err = parseDateOrDefault(dateTo, now); err != nil {
		return nil, err
	}

	listDeviceProgressStatusHistory, err := p.sDeviceProgressStatusHistoryRepo.FindByDate(ctx, dateFrom, dateTo)
	if err != nil {
		return nil, fmt.Errorf("p.CalcOEE: %w", err)
	}

	assignedWork, err := p.productionOrderStageDeviceRepo.GetAssignedByDate(ctx, dateFrom, dateTo)
	if err != nil {
		return nil, fmt.Errorf("p.CalcOEE: %w", err)
	}
	assignedWorkByDeviceID := processAssignedWorkByDeviceID(assignedWork)

	result := make(map[string]model.OEE)

	var lastHistory *repository.DeviceProgressStatusHistoryData

	for i := range listDeviceProgressStatusHistory {
		history := &listDeviceProgressStatusHistory[i]
		deviceID := history.DeviceID

		oee, exists := result[deviceID]
		if !exists {
			if lastHistory != nil {
				loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
				startOfDay := time.Date(lastHistory.CreatedAt.Year(), lastHistory.CreatedAt.Month(), lastHistory.CreatedAt.Day(), 7, 45, 0, 0, loc)

				deviceData := result[lastHistory.DeviceID]
				deviceData.ActualWorkingTime = lastHistory.CreatedAt.Sub(startOfDay).Milliseconds()
				result[lastHistory.DeviceID] = deviceData
			}

			oee = model.OEE{
				DowntimeStatistics: make(map[string]string),
				AssignedWork:       assignedWorkByDeviceID[deviceID],
			}

			for _, assigned := range assignedWorkByDeviceID[deviceID] {
				oee.TotalQuantity += assigned.Quantity
				if assigned.Settings != nil {
					if val, ok := assigned.Settings["san_pham_loi"].(int64); ok {
						oee.TotalDefective += val
					}
				}
				oee.AssignedWorkTime += assigned.EstimatedCompleteAt.Time.Sub(assigned.EstimatedStartAt.Time).Milliseconds()
			}
			result[deviceID] = oee
			lastHistory = history
			continue
		}
		if lastHistory == nil {
			continue
		}

		duration := history.CreatedAt.Sub(lastHistory.CreatedAt).Milliseconds()
		if history.ProcessStatus == enum.ProductionOrderStageDeviceStatusStart &&
			lastHistory.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed {
			oee.Downtime += duration
		}

		if history.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed ||
			history.ProcessStatus == enum.ProductionOrderStageDeviceStatusComplete {
			oee.JobRunningTime += duration

			if history.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed {
				oee.DowntimeStatistics[history.CreatedAt.Format(time.RFC3339)] = history.ErrorCode.String
			}
		}

		result[deviceID] = oee
		lastHistory = history
	}
	if lastHistory != nil {
		loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
		startOfDay := time.Date(lastHistory.CreatedAt.Year(), lastHistory.CreatedAt.Month(), lastHistory.CreatedAt.Day(), 7, 45, 0, 0, loc)

		deviceData := result[lastHistory.DeviceID]
		deviceData.ActualWorkingTime = lastHistory.CreatedAt.Sub(startOfDay).Milliseconds()
		result[lastHistory.DeviceID] = deviceData
	}
	return result, nil
}
func (p productionOrderStageDeviceService) CalcOEEByAssignedWork(ctx context.Context, dateFrom, dateTo string) (map[string]model.OEE, error) {
	now := time.Now().Format("2006-01-02")

	var err error
	if dateFrom, err = parseDateOrDefault(dateFrom, now); err != nil {
		return nil, err
	}
	if dateTo, err = parseDateOrDefault(dateTo, now); err != nil {
		return nil, err
	}

	listDeviceProgressStatusHistory, err := p.sDeviceProgressStatusHistoryRepo.FindByDate(ctx, dateFrom, dateTo)
	if err != nil {
		return nil, fmt.Errorf("p.CalcOEE: %w", err)
	}
	processDeviceProgressStatusHistory := processDeviceProgressStatusHistoryByProductionOrderStageDeviceID(listDeviceProgressStatusHistory)

	assignedWorks, err := p.productionOrderStageDeviceRepo.GetAssignedByDate(ctx, dateFrom, dateTo)
	if err != nil {
		return nil, fmt.Errorf("p.CalcOEE: %w", err)
	}

	result := make(map[string]model.OEE, len(assignedWorks))

	for _, assignedWork := range assignedWorks {
		defective := int64(0)
		if assignedWork.Settings != nil {
			if val, ok := assignedWork.Settings["san_pham_loi"].(int64); ok {
				defective = val
			}
		}

		oee := model.OEE{
			ActualWorkingTime:             28800,
			DowntimeStatistics:            make(map[string]string),
			AssignedWorkTime:              assignedWork.EstimatedCompleteAt.Time.Sub(assignedWork.EstimatedStartAt.Time).Milliseconds(),
			DeviceProgressStatusHistories: processDeviceProgressStatusHistory[assignedWork.ID],
			TotalQuantity:                 assignedWork.Quantity,
			TotalDefective:                defective,
		}

		histories := processDeviceProgressStatusHistory[assignedWork.ID]
		if len(histories) > 0 {
			lastHistory := histories[0]
			for _, history := range histories[1:] {
				duration := history.CreatedAt.Sub(lastHistory.CreatedAt).Milliseconds()

				switch {
				case history.ProcessStatus == enum.ProductionOrderStageDeviceStatusStart &&
					lastHistory.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed:
					oee.Downtime += duration

				case history.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed ||
					history.ProcessStatus == enum.ProductionOrderStageDeviceStatusComplete:
					oee.JobRunningTime += duration
					if history.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed {
						oee.DowntimeStatistics[history.CreatedAt.Format(time.RFC3339)] = history.ErrorCode.String
					}
				}

				lastHistory = history
			}
			loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
			startOfDay := time.Date(lastHistory.CreatedAt.Year(), lastHistory.CreatedAt.Month(), lastHistory.CreatedAt.Day(), 7, 45, 0, 0, loc)
			oee.ActualWorkingTime = lastHistory.CreatedAt.Sub(startOfDay).Milliseconds()
		}

		result[assignedWork.ID] = oee
	}

	return result, nil
}

func processAssignedWorkByDeviceID(datas []model.ProductionOrderStageDevice) map[string][]model.ProductionOrderStageDevice {
	result := make(map[string][]model.ProductionOrderStageDevice, len(datas))
	for i := range datas {
		result[datas[i].DeviceID] = append(result[datas[i].DeviceID], datas[i])
	}
	return result
}

func processDeviceProgressStatusHistoryByProductionOrderStageDeviceID(datas []repository.DeviceProgressStatusHistoryData) map[string][]model.DeviceProgressStatusHistory {
	result := make(map[string][]model.DeviceProgressStatusHistory, len(datas))
	for i := range datas {
		result[datas[i].ProductionOrderStageDeviceID] = append(result[datas[i].ProductionOrderStageDeviceID], *datas[i].DeviceProgressStatusHistory)
	}
	return result
}
