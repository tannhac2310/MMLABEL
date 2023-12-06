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
	IDs                  []string
	CustomerID           string
	ProductCode          string
	ProductName          string
	Name                 string
	EstimatedStartAtFrom time.Time //planned_production_date
	EstimatedStartAtTo   time.Time //planned_production_date
	Status               enum.ProductionOrderStatus
	Statuses             []enum.ProductionOrderStatus
	OrderStageStatus     enum.ProductionOrderStageStatus
	Responsible          []string
	StageIDs             []string
	Limit                int64
	Offset               int64
	Sort                 *Sort
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
		args = append(args, "%"+s.Name+"%")
		conds += fmt.Sprintf(" AND ( b.%[2]s ILIKE $%[1]d OR  b.%[3]s ILIKE $%[1]d OR b.%[4]s ILIKE $%[1]d OR b.%[5]s ILIKE $%[1]d OR b.%[6]s ILIKE $%[1]d) ",
			len(args), model.ProductionOrderFieldName, model.ProductionOrderFieldProductCode, model.ProductionOrderFieldProductName, model.ProductionOrderFieldCustomerID,
			model.ProductionOrderFieldSalesID,
		)
	}
	if s.CustomerID != "" {
		args = append(args, s.CustomerID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.ProductionOrderFieldCustomerID, len(args))
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
		conds += fmt.Sprintf(" AND b.%s < $%d", model.ProductionOrderFieldEstimatedStartAt, len(args))
	}

	if s.OrderStageStatus > 0 && len(s.StageIDs) > 0 {
		args = append(args, s.OrderStageStatus, s.StageIDs)
		joins += fmt.Sprintf(` INNER JOIN production_order_stages AS pos ON pos.production_order_id = b.id AND pos.deleted_at IS NULL`)
		conds += fmt.Sprintf(" AND pos.status = $%[1]d and pos.stage_id = ANY($%[2]d) ", len(args)-1, len(args))
	} else {
		if len(s.StageIDs) > 0 {
			args = append(args, s.StageIDs)
			joins += fmt.Sprintf(` INNER JOIN production_order_stages AS pos ON pos.production_order_id = b.id AND pos.deleted_at IS NULL`)
			conds += fmt.Sprintf("  AND pos.stage_id = ANY($%d)", len(args))
		}
		if s.OrderStageStatus > 0 {
			args = append(args, s.OrderStageStatus)
			joins += fmt.Sprintf(` INNER JOIN production_order_stages AS pos ON pos.production_order_id = b.id AND pos.deleted_at IS NULL`)
			conds += fmt.Sprintf(" AND pos.status = $%[1]d  ", len(args))
		}

	}

	if len(s.Responsible) > 0 {
		args = append(args, s.Responsible)

		joins += fmt.Sprintf(` INNER JOIN production_order_stages AS pos ON pos.production_order_id = b.id AND pos.deleted_at IS NULL 
			INNER JOIN production_order_stage_devices AS posd ON posd.production_order_stage_id = pos.id AND posd.deleted_at IS NULL`,
		)
		conds += fmt.Sprintf(" AND posd.%s && $%d", model.ProductionOrderStageDeviceFieldResponsible, len(args))

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
	return fmt.Sprintf(`SELECT b.%s
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
