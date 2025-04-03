package oee

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type Service interface {
	CalcOEEByDevice(ctx context.Context, opt *CalcOEEOpts) (map[string]model.OEE, map[string]string, error)
	CalcOEEByAssignedWork(ctx context.Context, opt *CalcOEEOpts) (map[string]model.OEE, int64, error)
}

type calcOEEService struct {
	oeeRepo repository.OEERepo
}

func NewService(
	oeeRepo repository.OEERepo,
) Service {
	return &calcOEEService{
		oeeRepo: oeeRepo,
	}

}

type CalcOEEOpts struct {
	ProductionOrderID            string
	ProductionOrderStageDeviceID string
	DateFrom                     string
	DateTo                       string
	DeviceID                     string
	Limit                        int64
	Offset                       int64
}

// parseDateOrDefault - Helper function to parse date or use default
func parseDateOrDefault(date string) (string, error) {
	defaultDate := time.Now().Format("2006-01-02")
	if date == "" {
		return defaultDate, nil
	}
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return "", fmt.Errorf("invalid date format: expected YYYY-MM-DD, got %s", date)
	}
	return date, nil
}

func (p calcOEEService) CalcOEEByDevice(ctx context.Context, opt *CalcOEEOpts) (map[string]model.OEE, map[string]string, error) {
	var err error
	// locUTC7, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	// if err != nil {
	// 	return nil, err
	// }
	//locUTC7 := time.FixedZone("UTC+7", 7*60*60)
	if opt.DateFrom, err = parseDateOrDefault(opt.DateFrom); err != nil {
		return nil, nil, err
	}
	if opt.DateTo, err = parseDateOrDefault(opt.DateTo); err != nil {
		return nil, nil, err
	}

	optRepo := repository.OEEOpts{
		ProductionOrderID:            opt.ProductionOrderID,
		ProductionOrderStageDeviceID: opt.ProductionOrderStageDeviceID,
		DateFrom:                     opt.DateFrom,
		DateTo:                       opt.DateTo,
		DeviceID:                     opt.DeviceID,
	}

	listDeviceProgressStatusHistory, err := p.oeeRepo.GetByDevice(ctx, optRepo)
	if err != nil {
		return nil, nil, fmt.Errorf("p.CalcOEE: %w", err)
	}

	assignedWork, _, err := p.oeeRepo.GetByAssigned(ctx, optRepo, -1, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("p.CalcOEE: %w", err)
	}
	assignedWorkByDeviceID, mapProductionOrderStageDevice, err := processAssignedWorkByDeviceID(ctx, assignedWork)
	if err != nil {
		return nil, nil, fmt.Errorf("p.CalcOEE: %w", err)
	}
	result := make(map[string]model.OEE)

	var lastHistory *repository.DeviceProgressStatusHistoryData
	for i := range listDeviceProgressStatusHistory {
		history := &listDeviceProgressStatusHistory[i]
		deviceID := history.DeviceID

		oee, exists := result[deviceID]
		if !exists {

			oee = model.OEE{
				DowntimeDetails:               make(map[string]int64),
				AssignedWork:                  assignedWorkByDeviceID[deviceID],
				DeviceProgressStatusHistories: make([]model.DeviceProgressStatusHistory, 0),
			}

			for _, assigned := range assignedWorkByDeviceID[deviceID] {
				oee.TotalQuantity += assigned.Quantity
				oee.TotalDefective += assigned.Defective
				oee.TotalAssignQuantity += assigned.AssignedQuantity
				oee.AssignedWorkTime += assigned.EstimatedCompleteAt.Sub(assigned.EstimatedStartAt).Milliseconds()
			}

			oee.DeviceProgressStatusHistories = append(oee.DeviceProgressStatusHistories, *history.DeviceProgressStatusHistory)
			result[deviceID] = oee
			lastHistory = history
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
			oee.DowntimeDetails[lastHistory.ErrorCode.String] += duration

		case (history.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed ||
			history.ProcessStatus == enum.ProductionOrderStageDeviceStatusComplete) &&
			lastHistory.ProcessStatus == enum.ProductionOrderStageDeviceStatusStart:
			oee.JobRunningTime += duration

			if history.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed {
				oee.DowntimeDetails[history.ErrorCode.String] += 0
			}
		}
		oee.DeviceProgressStatusHistories = append(oee.DeviceProgressStatusHistories, *history.DeviceProgressStatusHistory)
		result[deviceID] = oee

		lastHistory = history
	}

	return result, mapProductionOrderStageDevice, nil
}

func (o calcOEEService) CalcOEEByAssignedWork(ctx context.Context, opt *CalcOEEOpts) (map[string]model.OEE, int64, error) {
	var err error
	// locUTC7, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	// if err != nil {
	// 	return nil, 0, err
	// }
	//locUTC7 := time.FixedZone("UTC+7", 7*60*60)
	if opt.DateFrom, err = parseDateOrDefault(opt.DateFrom); err != nil {
		return nil, 0, err
	}
	if opt.DateTo, err = parseDateOrDefault(opt.DateTo); err != nil {
		return nil, 0, err
	}
	optRepo := repository.OEEOpts{
		ProductionOrderID:            opt.ProductionOrderID,
		ProductionOrderStageDeviceID: opt.ProductionOrderStageDeviceID,
		DateFrom:                     opt.DateFrom,
		DateTo:                       opt.DateTo,
		DeviceID:                     opt.DeviceID,
	}

	listDeviceProgressStatusHistory, err := o.oeeRepo.GetByDevice(ctx, optRepo)
	if err != nil {
		return nil, 0, fmt.Errorf("p.CalcOEE: %w", err)
	}
	processDeviceProgressStatusHistory := processDeviceProgressStatusHistoryByProductionOrderStageDeviceID(listDeviceProgressStatusHistory)

	assignedWorks, total, err := o.oeeRepo.GetByAssigned(ctx, optRepo, opt.Limit, opt.Offset)
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
						oee.DowntimeDetails[history.ErrorCode.String] += 0
					}
				}
				lastHistory = history
			}
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

func processAssignedWorkByDeviceID(ctx context.Context, datas []model.ProductionOrderStageDevice) (map[string][]model.AssignWorkOEE, map[string]string, error) {
	result := make(map[string][]model.AssignWorkOEE, len(datas))
	mapKeyProductOrderStage := make(map[string]string, 0)
	for i := range datas {
		var stageID string
		sqlQuery := `
			SELECT pos."stage_id" 
			FROM production_order_stages pos
			WHERE pos.id = $1;
		`
		if err := cockroach.Select(ctx, sqlQuery, datas[i].ProductionOrderStageID).ScanOne(&stageID); err != nil {
			return nil, nil, fmt.Errorf("p.CalcOEE: %w", err)
		}
		var defective int64 = 0
		if datas[i].Settings != nil {
			if val, ok := datas[i].Settings["san_pham_loi"].(int64); ok {
				defective = val
			}
		}
		result[datas[i].DeviceID] = append(result[datas[i].DeviceID], model.AssignWorkOEE{
			StageID:                stageID,
			ID:                     datas[i].ID,
			ProductionOrderStageID: datas[i].ProductionOrderStageID,
			EstimatedCompleteAt:    datas[i].EstimatedCompleteAt.Time,
			EstimatedStartAt:       datas[i].EstimatedStartAt.Time,
			Quantity:               datas[i].Quantity,
			Defective:              defective,
			AssignedQuantity:       datas[i].AssignedQuantity,
		})
		mapKeyProductOrderStage[datas[i].ID] = stageID
	}
	return result, mapKeyProductOrderStage, nil
}

func processDeviceProgressStatusHistoryByProductionOrderStageDeviceID(datas []repository.DeviceProgressStatusHistoryData) map[string][]model.DeviceProgressStatusHistory {
	result := make(map[string][]model.DeviceProgressStatusHistory, len(datas))
	for i := range datas {
		result[datas[i].ProductionOrderStageDeviceID] = append(result[datas[i].ProductionOrderStageDeviceID], *datas[i].DeviceProgressStatusHistory)
	}
	return result
}
