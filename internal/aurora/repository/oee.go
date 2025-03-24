package repository

import (
	"context"
	"fmt"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type OEERepo interface {
	GetByAssigned(ctx context.Context, opt OEEOpts, limit int64, offset int64) ([]model.ProductionOrderStageDevice, int64, error)
	GetByDevice(ctx context.Context, opt OEEOpts) ([]DeviceProgressStatusHistoryData, error)
}

type OEEOpts struct {
	ProductionOrderID            string
	ProductionOrderStageDeviceID string
	DateFrom                     string
	DateTo                       string
	DeviceID                     string
}

type oeeRepo struct {
}

func NewOEERepo() OEERepo {
	return &oeeRepo{}
}

func (o *oeeRepo) GetByAssigned(ctx context.Context, opt OEEOpts, limit, offset int64) ([]model.ProductionOrderStageDevice, int64, error) {
	sqlTable := `production_order_stage_devices posd`
	sqlCons := `posd.estimated_start_at::DATE BETWEEN $1 AND $2`
	args := []interface{}{opt.DateFrom, opt.DateTo}
	argIndex := 3

	if opt.DeviceID != "" {
		sqlCons += fmt.Sprintf(" AND posd.device_id = $%d", argIndex)
		args = append(args, opt.DeviceID)
		argIndex++
	}

	if opt.ProductionOrderStageDeviceID != "" {
		sqlCons += fmt.Sprintf(" AND posd.id = $%d", argIndex)
		args = append(args, opt.ProductionOrderStageDeviceID)
	}

	if opt.ProductionOrderID != "" {
		sqlTable += ` JOIN production_order_stage pos ON posd.production_order_stage_id = pos.id`
		sqlCons += fmt.Sprintf(" AND pos.production_order_id = $%d", argIndex)
		args = append(args, opt.ProductionOrderID)
		argIndex++
	}

	// Lấy tổng số bản ghi
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s`, sqlTable, sqlCons)
	var total int64
	if err := cockroach.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("cockroach.QueryRow (count): %w", err)
	}

	// Truy vấn dữ liệu với LIMIT và OFFSET
	sqlQuery := fmt.Sprintf(`
		SELECT posd.id, posd.production_order_stage_id, posd.device_id, posd.quantity, posd.settings, 
		       posd.estimated_start_at, posd.estimated_complete_at
		FROM %s 
		WHERE %s 
		ORDER BY posd.device_id, posd.estimated_start_at`, sqlTable, sqlCons)

	if limit > 0 {
		sqlQuery += fmt.Sprintf(" LIMIT %d", limit)
		if offset > 0 {
			sqlQuery += fmt.Sprintf(" OFFSET %d", offset)
		}
	}

	var result []model.ProductionOrderStageDevice
	if err := cockroach.Select(ctx, sqlQuery, args...).ScanAll(&result); err != nil {
		return nil, 0, fmt.Errorf("cockroach.Select: %w", err)
	}

	return result, total, nil
}

func (o *oeeRepo) GetByDevice(ctx context.Context, opt OEEOpts) ([]DeviceProgressStatusHistoryData, error) {
	var deviceProcessStatusHistoryData []DeviceProgressStatusHistoryData
	sqlTable := `device_progress_status_history dpsh`
	sqlCons := `dpsh.created_at::DATE BETWEEN $1 AND $2`
	args := []interface{}{opt.DateFrom, opt.DateTo}
	argIndex := 3

	if opt.DeviceID != "" {
		sqlCons += fmt.Sprintf(" AND dpsh.device_id = $%d", argIndex)
		args = append(args, opt.DeviceID)
		argIndex++
	}

	if opt.ProductionOrderStageDeviceID != "" {
		sqlCons += fmt.Sprintf(" AND dpsh.production_order_stage_device_id = $%d", argIndex)
		args = append(args, opt.ProductionOrderStageDeviceID)
		argIndex++
	}

	if opt.ProductionOrderID != "" {
		sqlTable += ` JOIN production_order_stage_device posd ON posd.id = dpsh.production_order_stage_device_id
		              JOIN production_order_stage pos ON posd.production_order_stage_id = pos.id`
		sqlCons += fmt.Sprintf(" AND pos.production_order_id = $%d", argIndex)
		args = append(args, opt.ProductionOrderID)
	}

	sqlQuery := fmt.Sprintf(`SELECT dpsh.* FROM %s WHERE %s ORDER BY dpsh.device_id, dpsh.created_at`, sqlTable, sqlCons)

	if err := cockroach.Select(ctx, sqlQuery, args...).ScanAll(&deviceProcessStatusHistoryData); err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return deviceProcessStatusHistoryData, nil
}
