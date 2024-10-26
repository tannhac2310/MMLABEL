package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type ProductRepo interface {
	Insert(ctx context.Context, e *model.Product) error
	Update(ctx context.Context, e *model.Product) error
	SoftDelete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*ProductData, error)
	Search(ctx context.Context, s *SearchProductOpts) ([]*ProductData, error)
	Count(ctx context.Context, s *SearchProductOpts) (*CountResult, error)
}

type sProductRepo struct {
}

func NewProductRepo() ProductRepo {
	return &sProductRepo{}
}

func (r *sProductRepo) Insert(ctx context.Context, e *model.Product) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sProductRepo) FindByID(ctx context.Context, id string) (*ProductData, error) {
	e := &ProductData{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("sProductRepo.cockroach.FindOne: %w", err)
	}

	return e, nil
}
func (r *sProductRepo) Update(ctx context.Context, e *model.Product) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *sProductRepo) SoftDelete(ctx context.Context, id string) error {
	sql := "UPDATE products SET deleted_at = NOW() WHERE id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("products cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*sProductRepo not found any records to delete")
	}

	return nil
}

// SearchProductOpts all params is options
type SearchProductOpts struct {
	IDs           []string
	Name          string
	CustomerID    string
	SaleID        string
	ProductPlanID string
	Limit         int64
	Offset        int64
	Sort          *Sort
}

func (s *SearchProductOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	joins += fmt.Sprintf("LEFT JOIN production_plan_products AS pp ON pp.product_id = b.id ")

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.ProductFieldID)
	}
	if s.Name != "" {
		args = append(args, "%"+s.Name+"%")
		conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.ProductFieldName, len(args))
	}
	if s.CustomerID != "" {
		args = append(args, s.CustomerID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.ProductFieldCustomerID, len(args))
	}
	if s.SaleID != "" {
		args = append(args, s.SaleID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.ProductFieldSaleID, len(args))
	}
	if s.ProductPlanID != "" {
		args = append(args, s.ProductPlanID)
		conds += fmt.Sprintf(" AND pp.production_plan_id = $%d", len(args))
	}

	b := &model.Product{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf("SELECT count(*) as cnt FROM %s AS b %s WHERE TRUE %s AND b.deleted_at IS NULL", b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.id DESC "
	if s.Sort != nil {
		order = fmt.Sprintf(" ORDER BY b.%s %s", s.Sort.By, s.Sort.Order)
	}
	return fmt.Sprintf("SELECT b.%s,pp.production_plan_id as production_plan_id   FROM %s AS b %s WHERE TRUE %s AND b.deleted_at IS NULL %s LIMIT %d OFFSET %d", strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type ProductData struct {
	*model.Product
	ProductionPlanID string `db:"production_plan_id"`
}

func (r *sProductRepo) Search(ctx context.Context, s *SearchProductOpts) ([]*ProductData, error) {
	Product := make([]*ProductData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&Product)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return Product, nil
}

func (r *sProductRepo) Count(ctx context.Context, s *SearchProductOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sProductRepo.Count: %w", err)
	}

	return countResult, nil
}
