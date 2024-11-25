package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type OrderItemRepo interface {
	Insert(ctx context.Context, e *model.OrderItem) error
	Update(ctx context.Context, e *model.OrderItem) error
	DeleteByOrderID(ctx context.Context, orderID string) error
	SortDeleteByOrderID(ctx context.Context, orderID string) error
	FindByID(ctx context.Context, id string) (*OrderItemData, error)
	Search(ctx context.Context, s *SearchOrderItemOpts) ([]*OrderItemData, error)
	Count(ctx context.Context, s *SearchOrderItemOpts) (*CountResult, error)
}

type sOrderItemRepo struct {
}

func NewOrderItemRepo() OrderItemRepo {
	return &sOrderItemRepo{}
}

func (r *sOrderItemRepo) Insert(ctx context.Context, e *model.OrderItem) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sOrderItemRepo) FindByID(ctx context.Context, id string) (*OrderItemData, error) {
	e := &OrderItemData{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("sOrderItemRepo.cockroach.FindOne: %w", err)
	}

	return e, nil
}
func (r *sOrderItemRepo) Update(ctx context.Context, e *model.OrderItem) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *sOrderItemRepo) DeleteByOrderID(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("sOrderItemRepo.DeleteByOrderID: id is required")
	}
	sql := "DELETE FROM order_items WHERE order_id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("order_items cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*sOrderItemRepo not found any records to delete")
	}

	return nil
}

func (r *sOrderItemRepo) SortDeleteByOrderID(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("sOrderItemRepo.SortDeleteByOrderID: id is required")
	}
	sql := "UPDATE order_items SET deleted_at = NOW() WHERE order_id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("order_items cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*sOrderItemRepo not found any records to delete")
	}

	return nil
}

// SearchOrderItemOpts all params is options
type SearchOrderItemOpts struct {
	IDs     []string
	OrderID string
	// todo add more search options
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchOrderItemOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.OrderItemFieldID)
	}
	// todo add more search options example:
	//if s.Name != "" {
	//	args = append(args, "%"+s.Name+"%")
	//	conds += fmt.Sprintf(" AND (b.%[2]s ILIKE $%[1]d OR b.%[3]s ILIKE $%[1]d)",
	//		len(args), model.OrderItemFieldName, model.OrderItemFieldCode)
	//}
	//if s.Code != "" {
	//	args = append(args, s.Code)
	//	conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.OrderItemFieldCode, len(args))
	//}
	if s.OrderID != "" {
		args = append(args, s.OrderID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.OrderItemFieldOrderID, len(args))
	}

	b := &model.OrderItem{}
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

type OrderItemData struct {
	*model.OrderItem
}

func (r *sOrderItemRepo) Search(ctx context.Context, s *SearchOrderItemOpts) ([]*OrderItemData, error) {
	OrderItem := make([]*OrderItemData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&OrderItem)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return OrderItem, nil
}

func (r *sOrderItemRepo) Count(ctx context.Context, s *SearchOrderItemOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sOrderItemRepo.Count: %w", err)
	}

	return countResult, nil
}
