package production_order_stage_device

import (
	"context"
	"fmt"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
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

	assignedWork, _, err := p.productionOrderStageDeviceRepo.GetAssignedByDate(ctx, dateFrom, dateTo, -1, -1)
	if err != nil {
		return nil, fmt.Errorf("p.CalcOEE: %w", err)
	}
	assignedWorkByDeviceID := processAssignedWorkByDeviceID(assignedWork)

	result := make(map[string]model.OEE)

	var lastHistory *repository.DeviceProgressStatusHistoryData
	var startOfDay time.Time
	for i := range listDeviceProgressStatusHistory {
		history := &listDeviceProgressStatusHistory[i]
		deviceID := history.DeviceID

		oee, exists := result[deviceID]
		if !exists {
			if lastHistory != nil {
				deviceData := result[lastHistory.DeviceID]
				deviceData.ActualWorkingTime = lastHistory.CreatedAt.Sub(startOfDay).Milliseconds()
				result[lastHistory.DeviceID] = deviceData
			}

			oee = model.OEE{
				DowntimeStatistics:            make(map[string]string),
				AssignedWork:                  assignedWorkByDeviceID[deviceID],
				DeviceProgressStatusHistories: make([]model.DeviceProgressStatusHistory, 0),
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
			oee.DeviceProgressStatusHistories = append(oee.DeviceProgressStatusHistories, *history.DeviceProgressStatusHistory)
			result[deviceID] = oee
			lastHistory = history
			startOfDay = history.CreatedAt
			continue
		}

		if lastHistory == nil {
			continue
		}

		duration := history.CreatedAt.Sub(lastHistory.CreatedAt).Milliseconds()
		switch {
		case history.ProcessStatus == enum.ProductionOrderStageDeviceStatusStart &&
			lastHistory.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed:
			oee.Downtime += duration

		case (history.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed ||
			history.ProcessStatus == enum.ProductionOrderStageDeviceStatusComplete) &&
			lastHistory.ProcessStatus == enum.ProductionOrderStageDeviceStatusStart:
			oee.JobRunningTime += duration

			if history.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed {
				oee.DowntimeStatistics[history.CreatedAt.Format(time.RFC3339)] = history.ErrorCode.String
			}
		}
		oee.DeviceProgressStatusHistories = append(oee.DeviceProgressStatusHistories, *history.DeviceProgressStatusHistory)
		result[deviceID] = oee

		lastHistory = history
	}
	return result, nil
}

func (p productionOrderStageDeviceService) CalcOEEByAssignedWork(ctx context.Context, dateFrom, dateTo string, limit int64, offset int64) (map[string]model.OEE, int64, error) {
	now := time.Now().Format("2006-01-02")

	var err error
	if dateFrom, err = parseDateOrDefault(dateFrom, now); err != nil {
		return nil, 0, err
	}
	if dateTo, err = parseDateOrDefault(dateTo, now); err != nil {
		return nil, 0, err
	}

	listDeviceProgressStatusHistory, err := p.sDeviceProgressStatusHistoryRepo.FindByDate(ctx, dateFrom, dateTo)
	if err != nil {
		return nil, 0, fmt.Errorf("p.CalcOEE: %w", err)
	}
	processDeviceProgressStatusHistory := processDeviceProgressStatusHistoryByProductionOrderStageDeviceID(listDeviceProgressStatusHistory)

	assignedWorks, total, err := p.productionOrderStageDeviceRepo.GetAssignedByDate(ctx, dateFrom, dateTo, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("p.CalcOEE: %w", err)
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
			DowntimeDetails:     make(map[string]int64),
			AssignedWorkTime:    assignedWork.EstimatedCompleteAt.Time.Sub(assignedWork.EstimatedStartAt.Time).Milliseconds(),
			TotalQuantity:       assignedWork.Quantity,
			TotalDefective:      defective,
			DeviceID:            assignedWork.DeviceID,
			TotalAssignQuantity: assignedWork.AssignedQuantity,
			Downtime:            0,
			JobRunningTime:      0,
		}
		histories := processDeviceProgressStatusHistory[assignedWork.ID]
		if len(histories) > 0 {
			lastHistory := histories[0]
			startOfDay := lastHistory.CreatedAt
			for _, history := range histories[1:] {
				duration := history.CreatedAt.Sub(lastHistory.CreatedAt).Milliseconds()

				switch {
				case history.ProcessStatus == enum.ProductionOrderStageDeviceStatusStart &&
					lastHistory.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed:
					oee.Downtime += duration
					oee.DowntimeDetails[lastHistory.ErrorCode.String] += duration

				case (history.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed ||
					history.ProcessStatus == enum.ProductionOrderStageDeviceStatusComplete) &&
					lastHistory.ProcessStatus == enum.ProductionOrderStageDeviceStatusStart:
					oee.JobRunningTime += duration
					if history.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed {
						oee.DowntimeDetails[lastHistory.ErrorCode.String] += 0
					}
				}
				lastHistory = history
			}
			oee.ActualWorkingTime = lastHistory.CreatedAt.Sub(startOfDay).Milliseconds()
		}

		var usernames []string
		sqlQuery := `SELECT DISTINCT u."name"
			FROM production_order_stage_responsible posr
			JOIN users u ON u.id = posr.user_id
			WHERE posr.po_stage_device_id = $1;
		`
		if err := cockroach.Select(ctx, sqlQuery, assignedWork.ID).ScanAll(&usernames); err != nil {
			return nil, 0, fmt.Errorf("p.CalcOEE: %w", err)
		}
		var productionOrderName string
		sqlQuery = `
			SELECT po."name" 
			FROM production_orders po 
			JOIN production_order_stages pos ON po.id = pos.production_order_id 
			WHERE pos.id = $1;
		`
		if err := cockroach.Select(ctx, sqlQuery, assignedWork.ProductionOrderStageID).ScanOne(&productionOrderName); err != nil {
			return nil, 0, fmt.Errorf("p.CalcOEE: %w", err)
		}
		oee.MachineOperator = usernames
		oee.ProductionOrderName = productionOrderName
		result[assignedWork.ID] = oee

	}
	return result, total, nil
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
