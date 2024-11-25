package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type OrderRepo interface {
	Insert(ctx context.Context, e *model.Order) error
	Update(ctx context.Context, e *model.Order) error
	SoftDelete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*OrderData, error)
	Search(ctx context.Context, s *SearchOrderOpts) ([]*OrderData, error)
	Count(ctx context.Context, s *SearchOrderOpts) (*CountResult, error)
	CntRows(ctx context.Context) (int64, error)
}

type sOrderRepo struct {
}

func NewOrderRepo() OrderRepo {
	return &sOrderRepo{}
}

func (r *sOrderRepo) Insert(ctx context.Context, e *model.Order) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sOrderRepo) FindByID(ctx context.Context, id string) (*OrderData, error) {
	e := &OrderData{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("sOrderRepo.cockroach.FindOne: %w", err)
	}

	return e, nil
}
func (r *sOrderRepo) Update(ctx context.Context, e *model.Order) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *sOrderRepo) SoftDelete(ctx context.Context, id string) error {
	sql := "UPDATE orders SET deleted_at = NOW() WHERE id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("orders cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*sOrderRepo not found any records to delete")
	}

	return nil
}

// SearchOrderOpts all params is options
type SearchOrderOpts struct {
	IDs []string
	// todo add more search options
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchOrderOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.OrderFieldID)
	}
	// todo add more search options example:
	//if s.Name != "" {
	//	args = append(args, "%"+s.Name+"%")
	//	conds += fmt.Sprintf(" AND (b.%[2]s ILIKE $%[1]d OR b.%[3]s ILIKE $%[1]d)",
	//		len(args), model.OrderFieldName, model.OrderFieldCode)
	//}
	//if s.Code != "" {
	//	args = append(args, s.Code)
	//	conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.OrderFieldCode, len(args))
	//}

	b := &model.Order{}
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

type OrderData struct {
	*model.Order
}

func (r *sOrderRepo) Search(ctx context.Context, s *SearchOrderOpts) ([]*OrderData, error) {
	Order := make([]*OrderData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&Order)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return Order, nil
}

func (r *sOrderRepo) Count(ctx context.Context, s *SearchOrderOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sOrderRepo.Count: %w", err)
	}

	return countResult, nil
}

func (r *sOrderRepo) CntRows(ctx context.Context) (int64, error) {
	sql := "SELECT count(*) FROM orders"
	var count int64
	err := cockroach.Select(ctx, sql).ScanOne(&count)
	if err != nil {
		return 0, fmt.Errorf("count rows: %w", err)
	}

	return count, nil
}
