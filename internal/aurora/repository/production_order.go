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
	Count(ctx context.Context, s *SearchProductionOrdersOpts) (*CountResult, error)
}

type productionOrdersRepo struct {
}

func NewProductionOrderRepo() ProductionOrderRepo {
	return &productionOrdersRepo{}
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
	IDs         []string
	CustomerID  string
	ProductCode string
	ProductName string
	Name        string
	Status      enum.ProductionOrderStatus
	Limit       int64
	Offset      int64
	Sort        *Sort
}

func (s *SearchProductionOrdersOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.ProductionOrderFieldID)
	}
	if s.Name != "" {
		args = append(args, "%"+s.Name+"%")
		conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.ProductionOrderFieldName, len(args))
	}
	if s.CustomerID != "" {
		args = append(args, s.CustomerID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.ProductionOrderFieldCustomerID, len(args))
	}

	b := &model.ProductionOrder{}
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
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&message)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return message, nil
}

func (r *productionOrdersRepo) Count(ctx context.Context, s *SearchProductionOrdersOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("chat.Count: %w", err)
	}

	return countResult, nil
}
