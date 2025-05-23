package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type ProductionOrderRepo interface {
	Insert(ctx context.Context, e *model.ProductionOrder) error
	Update(ctx context.Context, e *model.ProductionOrder) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchProductionOrdersOpts) ([]*ProductionOrderData, error)
	Analysis(ctx context.Context, s *SearchProductionOrdersOpts) ([]*Analysis, error)
	Count(ctx context.Context, s *SearchProductionOrdersOpts) (*CountResult, error)
	CountByCreatedDate(ctx context.Context, from, to time.Time) (int64, error)
	CountByCode(ctx context.Context, code string) (int64, error)
}

type productionOrdersRepo struct {
}

func NewProductionOrderRepo() ProductionOrderRepo {
	return &productionOrdersRepo{}
}

type Analysis struct {
	Status enum.ProductionOrderStatus `db:"status"`
	Count  int64                      `db:"count"`
}

func (r *productionOrdersRepo) Analysis(ctx context.Context, s *SearchProductionOrdersOpts) ([]*Analysis, error) {
	data := make([]*Analysis, 0)
	sql, args := s.buildQuery(false, true)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&data)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Analysis: %w", err)
	}

	return data, nil
}

func (r *productionOrdersRepo) Insert(ctx context.Context, e *model.ProductionOrder) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *productionOrdersRepo) Update(ctx context.Context, e *model.ProductionOrder) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *productionOrdersRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE production_orders
		SET deleted_at = NOW()
		WHERE id = $1`

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found any records to delete")
	}

	return nil
}

// SearchProductionOrdersOpts all params is options
type SearchProductionOrdersOpts struct {
	IDs                             []string
	ProductionPlanIDs               []string
	CustomerID                      string
	ProductCode                     string
	ProductName                     string
	Name                            string
	EstimatedStartAtFrom            time.Time //planned_production_date
	EstimatedStartAtTo              time.Time //planned_production_date
	EstimatedCompletedFrom          time.Time
	EstimatedCompletedTo            time.Time
	Status                          enum.ProductionOrderStatus
	Statuses                        []enum.ProductionOrderStatus
	OrderStageStatus                enum.ProductionOrderStageStatus
	OrderStageEstimatedStartFrom    time.Time
	OrderStageEstimatedStartTo      time.Time
	OrderStageEstimatedCompleteFrom time.Time
	OrderStageEstimatedCompleteTo   time.Time
	Responsible                     []string
	StageIDs                        []string
	StageInLine                     string
	UserID                          string
	DeviceID                        string
	Limit                           int64
	Offset                          int64
	Sort                            *Sort
}

func (s *SearchProductionOrdersOpts) buildQuery(isCount bool, isAnalysis bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.ProductionOrderFieldID)
	}

	if s.Name != "" {
		args = append(args, "%"+strings.Trim(s.Name, " ")+"%")
		joins += " JOIN custom_fields AS cf ON cf.entity_id = b.id "
		conds += fmt.Sprintf(" AND ( b.%[2]s ILIKE $%[1]d OR b.%[3]s ILIKE $%[1]d OR b.%[4]s ILIKE $%[1]d OR b.%[5]s ILIKE $%[1]d OR b.%[6]s ILIKE $%[1]d OR b.%[7]s ILIKE $%[1]d OR cf.value ILIKE $%[1]d )",
			len(args), model.ProductionOrderFieldName, model.ProductionOrderFieldProductCode, model.ProductionOrderFieldProductName, model.ProductionOrderFieldCustomerID,
			model.ProductionOrderFieldSalesID, model.ProductionOrderFieldID,
		)
	}

	if s.CustomerID != "" {
		args = append(args, s.CustomerID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.ProductionOrderFieldCustomerID, len(args))
	}

	if s.Status > 0 && !isAnalysis { // neu isAnalysis = true thi khong can check status
		args = append(args, s.Status)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.ProductionOrderFieldStatus, len(args))
	}

	if len(s.ProductionPlanIDs) > 0 {
		args = append(args, s.ProductionPlanIDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($%d)", model.ProductionOrderFieldProductionPlanID, len(args))
	}

	if len(s.Statuses) > 0 {
		args = append(args, s.Statuses)
		conds += fmt.Sprintf(" AND b.%s = ANY($%d)", model.ProductionOrderFieldStatus, len(args))
	}

	if !s.EstimatedStartAtFrom.IsZero() {
		args = append(args, s.EstimatedStartAtFrom)
		conds += fmt.Sprintf(" AND b.%s >= $%d", model.ProductionOrderFieldEstimatedStartAt, len(args))
	}
	if !s.EstimatedStartAtTo.IsZero() {
		args = append(args, s.EstimatedStartAtTo)
		conds += fmt.Sprintf(" AND b.%s <= $%d", model.ProductionOrderFieldEstimatedStartAt, len(args))
	}

	if !s.EstimatedCompletedFrom.IsZero() {
		args = append(args, s.EstimatedCompletedFrom)
		conds += fmt.Sprintf(" AND b.%s >= $%d", model.ProductionOrderFieldEstimatedCompleteAt, len(args))
	}
	if !s.EstimatedCompletedTo.IsZero() {
		args = append(args, s.EstimatedCompletedTo)
		conds += fmt.Sprintf(" AND b.%s <= $%d", model.ProductionOrderFieldEstimatedCompleteAt, len(args))
	}

	if s.OrderStageStatus > 0 && len(s.StageIDs) > 0 {
		args = append(args, s.OrderStageStatus, s.StageIDs)
		conds += fmt.Sprintf(` AND EXISTS (SELECT 1 FROM production_order_stages AS pos WHERE pos.production_order_id = b.id AND pos.deleted_at IS NULL AND pos.status = $%[1]d AND pos.stage_id = ANY($%[2]d))`, len(args)-1, len(args))

	} else {
		if len(s.StageIDs) > 0 {
			args = append(args, s.StageIDs)
			conds += fmt.Sprintf(` AND EXISTS (SELECT 1 FROM production_order_stages AS pos WHERE pos.production_order_id = b.id AND pos.deleted_at IS NULL AND pos.stage_id = ANY($%[1]d))`, len(args))
		}
		if s.OrderStageStatus > 0 {
			args = append(args, s.OrderStageStatus)
			conds += fmt.Sprintf(` AND EXISTS (SELECT 1 FROM production_order_stages AS pos WHERE pos.production_order_id = b.id AND pos.deleted_at IS NULL AND pos.status = $%[1]d)`, len(args))
		}
	}
	if !s.OrderStageEstimatedStartFrom.IsZero() {
		args = append(args, s.OrderStageEstimatedStartFrom)
		conds += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM production_order_stages AS pos WHERE pos.production_order_id = b.id AND pos.deleted_at IS NULL AND pos.estimated_start_at >= $%d)", len(args))
	}
	if !s.OrderStageEstimatedStartTo.IsZero() {
		args = append(args, s.OrderStageEstimatedStartTo)
		conds += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM production_order_stages AS pos WHERE pos.production_order_id = b.id AND pos.deleted_at IS NULL AND pos.estimated_start_at < $%d)", len(args))
	}
	if !s.OrderStageEstimatedCompleteFrom.IsZero() {
		args = append(args, s.OrderStageEstimatedCompleteFrom)
		conds += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM production_order_stages AS pos WHERE pos.production_order_id = b.id AND pos.deleted_at IS NULL AND pos.estimated_complete_at >= $%d)", len(args))
	}
	if !s.OrderStageEstimatedCompleteTo.IsZero() {
		args = append(args, s.OrderStageEstimatedCompleteTo)
		conds += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM production_order_stages AS pos WHERE pos.production_order_id = b.id AND pos.deleted_at IS NULL AND pos.estimated_complete_at < $%d)", len(args))
	}

	// filter by device_id in production_order_stage_devices
	if s.DeviceID != "" {
		args = append(args, s.DeviceID)
		conds += fmt.Sprintf(` AND EXISTS (SELECT 1 FROM production_order_stages AS pos WHERE pos.production_order_id = b.id AND pos.deleted_at IS NULL
						AND EXISTS (SELECT 1 FROM production_order_stage_devices AS posd WHERE posd.production_order_stage_id = pos.id AND posd.deleted_at IS NULL AND posd.%[1]s = $%[2]d))`,
			model.ProductionOrderStageDeviceFieldDeviceID, len(args))
	}

	if s.StageInLine != "" && len(s.Responsible) > 0 {
		args = append(args, s.StageInLine, s.Responsible)
		// and pos.status = 3 đang thực hiện
		// posd.stage_id = $%[1]d AND posd.responsible && $%[2]d : trong công đoạn đang thực hiện X, và user có trong mảng responsible
		conds += fmt.Sprintf(` AND EXISTS (SELECT 1 FROM production_order_stages AS pos WHERE pos.production_order_id = b.id AND pos.deleted_at IS NULL and pos.status = 3
						AND EXISTS (SELECT 1 FROM production_order_stage_devices AS posd WHERE posd.production_order_stage_id = pos.id AND posd.deleted_at IS NULL AND pos.stage_id = $%[1]d AND posd.responsible && $%[2]d))`,
			len(args)-1, len(args))
	} else {
		// StageInLine công đoạn đang thực hiện pos.status = enum.ProductionOrderStageStatusProductionStart = 3
		if s.StageInLine != "" {
			args = append(args, s.StageInLine)
			conds += fmt.Sprintf(` AND EXISTS (SELECT 1 FROM production_order_stages AS pos WHERE
					pos.production_order_id = b.id AND pos.deleted_at IS NULL AND pos.stage_id = $%[1]d and pos.status = 3)`, len(args)) // 3: production_start
		}

		if len(s.Responsible) > 0 {
			args = append(args, s.Responsible)
			conds += fmt.Sprintf(` AND EXISTS (SELECT 1 FROM production_order_stages AS pos WHERE pos.production_order_id = b.id AND pos.deleted_at IS NULL
						AND EXISTS (SELECT 1 FROM production_order_stage_devices AS posd WHERE posd.production_order_stage_id = pos.id AND posd.deleted_at IS NULL AND posd.%[1]s && $%[2]d))`,
				model.ProductionOrderStageDeviceFieldResponsible, len(args))
		}
	}

	b := &model.ProductionOrder{}
	if isAnalysis {
		return fmt.Sprintf(`SELECT b.status, count(*) as count
		FROM %s AS b %s
		WHERE TRUE %s AND b.deleted_at IS NULL GROUP BY b.status`, b.TableName(), joins, conds), args
	}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf(`SELECT count(*) as cnt
		FROM %s AS b %s
		WHERE TRUE %s AND b.deleted_at IS NULL`, b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.id DESC "
	if s.Sort != nil {
		order = fmt.Sprintf(" ORDER BY b.%s %s", s.Sort.By, s.Sort.Order)
	}
	return fmt.Sprintf(`SELECT DISTINCT b.%s
		FROM %s AS b %s
		WHERE TRUE %s AND b.deleted_at IS NULL
		%s
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type ProductionOrderData struct {
	*model.ProductionOrder
}

func (r *productionOrdersRepo) Search(ctx context.Context, s *SearchProductionOrdersOpts) ([]*ProductionOrderData, error) {
	message := make([]*ProductionOrderData, 0)
	sql, args := s.buildQuery(false, false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&message)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return message, nil
}

func (r *productionOrdersRepo) Count(ctx context.Context, s *SearchProductionOrdersOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true, false)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("chat.Count: %w", err)
	}

	return countResult, nil
}

func (r *productionOrdersRepo) CountByCreatedDate(ctx context.Context, from, to time.Time) (int64, error) {
	sql := `SELECT count(*) FROM production_orders WHERE created_at >= $1 AND created_at <= $2`
	var count int64
	err := cockroach.Select(ctx, sql, from, to).ScanOne(&count)
	if err != nil {
		return 0, fmt.Errorf("cockroach.Select: %w", err)
	}

	return count, nil
}

func (r *productionOrdersRepo) CountByCode(ctx context.Context, code string) (int64, error) {
	sql := `SELECT count(*) FROM production_orders WHERE product_code like $1`
	var count int64
	err := cockroach.Select(ctx, sql, "%"+code+"%").ScanOne(&count)
	if err != nil {
		return 0, fmt.Errorf("cockroach.Select: %w", err)
	}

	return count, nil
}
