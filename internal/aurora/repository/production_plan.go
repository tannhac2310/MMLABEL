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

type ProductionPlanRepo interface {
	Insert(ctx context.Context, e *model.ProductionPlan) error
	Update(ctx context.Context, e *model.ProductionPlan) error
	SoftDelete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*ProductionPlanData, error)
	Search(ctx context.Context, s *SearchProductionPlanOpts) ([]*ProductionPlanData, error)
	Count(ctx context.Context, s *SearchProductionPlanOpts) (*CountResult, error)
}

type sProductionPlanRepo struct {
}

func NewProductionPlanRepo() ProductionPlanRepo {
	return &sProductionPlanRepo{}
}

func (r *sProductionPlanRepo) Insert(ctx context.Context, e *model.ProductionPlan) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sProductionPlanRepo) FindByID(ctx context.Context, id string) (*ProductionPlanData, error) {
	e := &ProductionPlanData{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("sProductionPlanRepo.cockroach.FindOne: %w", err)
	}

	return e, nil
}
func (r *sProductionPlanRepo) Update(ctx context.Context, e *model.ProductionPlan) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *sProductionPlanRepo) SoftDelete(ctx context.Context, id string) error {
	sql := "UPDATE production_plans SET deleted_at = NOW() WHERE id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("production_plans cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*sProductionPlanRepo not found any records to delete")
	}

	return nil
}

// SearchProductionPlanOpts all params is options
type SearchProductionPlanOpts struct {
	IDs        []string
	CustomerID string
	Name       string
	Statuses   []enum.ProductionOrderStatus
	UserID     string
	Limit      int64
	Offset     int64
	Sort       *Sort
}

func (s *SearchProductionPlanOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.ProductionPlanFieldID)
	}
	if s.Name != "" {
		args = append(args, "%"+s.Name+"%")
		conds += fmt.Sprintf(" AND b.%[2]s ILIKE $%[1]d", len(args), model.ProductionPlanFieldName)
	}
	if len(s.Statuses) > 0 {
		args = append(args, s.Statuses)
		conds += fmt.Sprintf(" AND b.%s = ANY($%d)", model.ProductionPlanFieldStatus, len(args))
	}

	b := &model.ProductionPlan{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf("SELECT count(*) as cnt FROM %s AS b %s WHERE TRUE %s AND b.deleted_at IS NULL", b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.id DESC "
	if s.Sort != nil {
		order = fmt.Sprintf(" ORDER BY b.%s %s", s.Sort.By, s.Sort.Order)
	}
	return fmt.Sprintf("SELECT b.%s FROM %s AS b %s WHERE TRUE %s AND b.deleted_at IS NULL %s LIMIT %d OFFSET %d", strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type ProductionPlanData struct {
	*model.ProductionPlan
}

func (r *sProductionPlanRepo) Search(ctx context.Context, s *SearchProductionPlanOpts) ([]*ProductionPlanData, error) {
	ProductionPlan := make([]*ProductionPlanData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&ProductionPlan)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return ProductionPlan, nil
}

func (r *sProductionPlanRepo) Count(ctx context.Context, s *SearchProductionPlanOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sProductionPlanRepo.Count: %w", err)
	}

	return countResult, nil
}
