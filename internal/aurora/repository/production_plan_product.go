package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type ProductionPlanProductRepo interface {
	Insert(ctx context.Context, e *model.ProductionPlanProduct) error
	Update(ctx context.Context, e *model.ProductionPlanProduct) error
	SoftDelete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*ProductionPlanProductData, error)
	Search(ctx context.Context, s *SearchProductionPlanProductOpts) ([]*ProductionPlanProductData, error)
	Count(ctx context.Context, s *SearchProductionPlanProductOpts) (*CountResult, error)
}

type sProductionPlanProductRepo struct {
}

func NewProductionPlanProductRepo() ProductionPlanProductRepo {
	return &sProductionPlanProductRepo{}
}

func (r *sProductionPlanProductRepo) Insert(ctx context.Context, e *model.ProductionPlanProduct) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sProductionPlanProductRepo) FindByID(ctx context.Context, id string) (*ProductionPlanProductData, error) {
	e := &ProductionPlanProductData{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("sProductionPlanProductRepo.cockroach.FindOne: %w", err)
	}

	return e, nil
}
func (r *sProductionPlanProductRepo) Update(ctx context.Context, e *model.ProductionPlanProduct) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *sProductionPlanProductRepo) SoftDelete(ctx context.Context, id string) error {
	sql := "UPDATE production_plan_products SET deleted_at = NOW() WHERE id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("production_plan_products cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*sProductionPlanProductRepo not found any records to delete")
	}

	return nil
}

// SearchProductionPlanProductOpts all params is options
type SearchProductionPlanProductOpts struct {
	IDs []string
	// todo add more search options
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchProductionPlanProductOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.ProductionPlanProductFieldID)
	}
	// todo add more search options example:
	//if s.Name != "" {
	//	args = append(args, "%"+s.Name+"%")
	//	conds += fmt.Sprintf(" AND (b.%[2]s ILIKE $%[1]d OR b.%[3]s ILIKE $%[1]d)",
	//		len(args), model.ProductionPlanProductFieldName, model.ProductionPlanProductFieldCode)
	//}
	//if s.Code != "" {
	//	args = append(args, s.Code)
	//	conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.ProductionPlanProductFieldCode, len(args))
	//}

	b := &model.ProductionPlanProduct{}
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

type ProductionPlanProductData struct {
	*model.ProductionPlanProduct
}

func (r *sProductionPlanProductRepo) Search(ctx context.Context, s *SearchProductionPlanProductOpts) ([]*ProductionPlanProductData, error) {
	ProductionPlanProduct := make([]*ProductionPlanProductData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&ProductionPlanProduct)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return ProductionPlanProduct, nil
}

func (r *sProductionPlanProductRepo) Count(ctx context.Context, s *SearchProductionPlanProductOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sProductionPlanProductRepo.Count: %w", err)
	}

	return countResult, nil
}
