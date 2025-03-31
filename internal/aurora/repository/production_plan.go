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
	CountRows(ctx context.Context) int64
	Summary(ctx context.Context, s *SummaryProductionPlanOpts) ([]*SummaryProductionPlanData, error)
}

type sProductionPlanRepo struct {
}

func NewProductionPlanRepo() ProductionPlanRepo {
	return &sProductionPlanRepo{}
}
func (r *sProductionPlanRepo) CountRows(ctx context.Context) int64 {
	countResult := &CountResult{}
	sql := "SELECT count(*) as cnt FROM production_plans AS b"

	err := cockroach.Select(ctx, sql).ScanOne(countResult)
	if err != nil {
		return 0
	}

	return countResult.Count

}
func (r *sProductionPlanRepo) Insert(ctx context.Context, e *model.ProductionPlan) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sProductionPlanRepo) FindByID(ctx context.Context, id string) (*ProductionPlanData, error) {
	e := &ProductionPlanData{
		ProductionPlan: &model.ProductionPlan{},
	}
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

type SummaryProductionPlanOpts struct {
	StartDate time.Time
	EndDate   time.Time
}

func (s *SummaryProductionPlanOpts) buildQuery() (string, []interface{}) {
	var args []interface{}
	conds := ""
	if !s.StartDate.IsZero() && !s.EndDate.IsZero() {
		args = append(args, s.StartDate, s.EndDate)
		conds += fmt.Sprintf(" AND b.%s BETWEEN $%d AND $%d", model.ProductionPlanFieldCreatedAt, len(args)-1, len(args))
	} else if !s.StartDate.IsZero() {
		args = append(args, s.StartDate)
		conds += fmt.Sprintf(" AND b.%s >= $%d", model.ProductionPlanFieldCreatedAt, len(args))
	} else if !s.EndDate.IsZero() {
		args = append(args, s.EndDate)
		conds += fmt.Sprintf(" AND b.%s <= $%d", model.ProductionPlanFieldCreatedAt, len(args))
	}

	b := &model.ProductionPlan{}
	return fmt.Sprintf("SELECT b.current_stage, b.status, count(1) AS count FROM %s AS b WHERE TRUE %s AND b.deleted_at IS NULL group by b.current_stage, b.status", b.TableName(), conds), args
}

type SummaryProductionPlanData struct {
	Stage  enum.ProductionPlanStage  `db:"current_stage"`
	Status enum.ProductionPlanStatus `db:"status"`
	Count  int64                     `db:"count"`
}

func (r *sProductionPlanRepo) Summary(ctx context.Context, s *SummaryProductionPlanOpts) ([]*SummaryProductionPlanData, error) {
	SummaryProductionPlan := make([]*SummaryProductionPlanData, 0)
	sql, args := s.buildQuery()
	err := cockroach.Select(ctx, sql, args...).ScanAll(&SummaryProductionPlan)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return SummaryProductionPlan, nil
}

// SearchProductionPlanOpts all params is options
type SearchProductionPlanOpts struct {
	IDs []string
	//CustomerID  string
	Name        string
	Search      string
	ProductName string
	ProductCode string
	Statuses    []enum.ProductionPlanStatus
	//UserID      string // TODO add later
	Stage  enum.ProductionPlanStage
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchProductionPlanOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := " LEFT JOIN users AS cu ON b.created_by = cu.id LEFT JOIN users AS uu ON b.updated_by = uu.id "

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.ProductionPlanFieldID)
	}
	if s.Name != "" {
		args = append(args, "%"+s.Name+"%")
		conds += fmt.Sprintf(" AND (b.%[2]s ILIKE $%[1]d OR b.id LIKE 'in-' || $%[1]d)", len(args), model.ProductionPlanFieldName)
	}

	if s.Search != "" {
		args = append(args, "%"+s.Search+"%")
		conds += fmt.Sprintf(" AND (b.%[2]s ILIKE $%[1]d OR id = 'in+' || $%[1]d)", len(args), model.ProductionPlanFieldSearchContent)
	}

	if s.ProductName != "" {
		args = append(args, "%"+s.ProductName+"%")
		conds += fmt.Sprintf(" AND b.%[2]s ILIKE $%[1]d", len(args), model.ProductionPlanFieldProductName)
	}
	if s.ProductCode != "" {
		args = append(args, s.ProductCode)
		conds += fmt.Sprintf(" AND b.%[2]s ILIKE $%[1]d", len(args), model.ProductionPlanFieldProductCode)
	}
	if len(s.Statuses) > 0 {
		args = append(args, s.Statuses)
		conds += fmt.Sprintf(" AND b.%s = ANY($%d)", model.ProductionPlanFieldStatus, len(args))
	}
	if s.Stage > 0 {
		args = append(args, s.Stage)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.ProductionPlanFieldCurrentStage, len(args))
	}
	// TODO need to join production_plan_products products

	//if s.CustomerID != "" {
	//	// customer id in products
	//	args = append(args, s.CustomerID)
	//	conds += fmt.Sprintf(" AND b.%s = $%d", model.ProductionPlanFieldCustomerID, len(args))
	//}

	b := &model.ProductionPlan{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf("SELECT count(*) as cnt FROM %s AS b %s WHERE TRUE %s AND b.deleted_at IS NULL", b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.created_at DESC "
	if s.Sort != nil && s.Sort.By != "" && s.Sort.Order != "" {
		order = fmt.Sprintf(" ORDER BY b.%s %s", s.Sort.By, s.Sort.Order)
	}
	return fmt.Sprintf("SELECT b.%s, cu.name as created_by_name, uu.name as updated_by_name FROM %s AS b %s WHERE TRUE %s AND b.deleted_at IS NULL %s LIMIT %d OFFSET %d", strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type ProductionPlanData struct {
	*model.ProductionPlan
	CreatedByName string `db:"created_by_name"`
	UpdatedByName string `db:"updated_by_name"`
}

func (r *sProductionPlanRepo) Search(ctx context.Context, s *SearchProductionPlanOpts) ([]*ProductionPlanData, error) {
	ProductionPlan := make([]*ProductionPlanData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&ProductionPlan)
	if err != nil {
		return nil, fmt.Errorf("production plan search: %w", err)
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
