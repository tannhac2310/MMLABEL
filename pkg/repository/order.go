package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

type OrderRepo interface {
	Insert(ctx context.Context, e *model.Order) error
	Update(ctx context.Context, e *model.Order) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchOrdersOpts) ([]*OrderData, error)
	Count(ctx context.Context, s *SearchOrdersOpts) (*CountResult, error)
	ReportOrderByDate(ctx context.Context) ([]*OrderByDate, error)
}
type OrderByDate struct {
	Date   time.Time        `db:"date"`
	Status enum.OrderStatus `db:"status"`
	Total  float64          `db:"total"`
}

type OrderData struct {
	*model.Order
	CreatedByName string
	UpdatedByName string
}

type orderRepo struct {
}

func (r *orderRepo) ReportOrderByDate(ctx context.Context) ([]*OrderByDate, error) {
	sql := `
select sum(price - discount) total,
       status,
       experimental_strftime(created_at, '%Y-%m-%d')
           as                date
from orders
group by date, status
`
	data := make([]*OrderByDate, 0)
	err := cockroach.Select(ctx, sql).ScanAll(&data)
	return data, err
}

func NewOrderRepo() OrderRepo {
	return &orderRepo{}
}

func (r *orderRepo) Insert(ctx context.Context, e *model.Order) error {

	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *orderRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE orders
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
func (r *orderRepo) Update(ctx context.Context, e *model.Order) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

// SearchOrdersOpts all params is options
type SearchOrdersOpts struct {
	IDs             []string
	StudentID       string
	StudentFullName string
	StudentPhone    string
	StudentEmail    string
	InvoiceID       string
	PaymentMethod   enum.PaymentMethod
	SearchString    string
	Sort            *Sort
	Limit           int64
	Offset          int64
}

func (s *SearchOrdersOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND o.%s = ANY($1)", model.OrderFieldID)
	}

	if s.SearchString != "" {
		args = append(args, "%"+s.SearchString+"%")
		conds += fmt.Sprintf(" AND o.%s ILIKE $%d", model.OrderFieldStudentFullName, len(args))

		args = append(args, "%"+s.SearchString+"%")
		conds += fmt.Sprintf(" AND o.%s ILIKE $%d", model.OrderFieldAddress, len(args))

		args = append(args, "%"+s.SearchString+"%")
		conds += fmt.Sprintf(" AND o.%s ILIKE $%d", model.OrderFieldNote, len(args))

		args = append(args, "%"+s.SearchString+"%")
		conds += fmt.Sprintf(" AND o.%s ILIKE $%d", model.OrderFieldStudentEmail, len(args))

		args = append(args, "%"+s.SearchString+"%")
		conds += fmt.Sprintf(" AND o.%s ILIKE $%d", model.OrderFieldStudentPhone, len(args))
	}

	if s.StudentFullName != "" {
		args = append(args, s.StudentFullName)
		conds += fmt.Sprintf(" AND o.%s = $%d", model.OrderFieldStudentFullName, len(args))
	}

	if s.StudentID != "" {
		args = append(args, s.StudentID)
		conds += fmt.Sprintf(" AND o.%s = $%d", model.OrderFieldStudentID, len(args))
	}

	if s.StudentEmail != "" {
		args = append(args, s.StudentEmail)
		conds += fmt.Sprintf(" AND o.%s = $%d", model.OrderFieldStudentEmail, len(args))
	}

	if s.InvoiceID != "" {
		args = append(args, s.InvoiceID)
		conds += fmt.Sprintf(" AND o.%s = $%d", model.OrderFieldInvoiceID, len(args))
	}
	if s.PaymentMethod > 0 {
		args = append(args, s.PaymentMethod)
		conds += fmt.Sprintf(" AND o.%s = $%d", model.OrderFieldPaymentMethod, len(args))
	}
	c := &model.Order{}
	fields, _ := c.FieldMap()
	if isCount {
		return fmt.Sprintf(`SELECT count(*) as cnt
		FROM %s AS o %s
		WHERE TRUE %s AND o.deleted_at IS NULL`, c.TableName(), joins, conds), args
	}

	order := " ORDER BY o.id DESC "
	if s.Sort != nil {
		order = fmt.Sprintf(" ORDER BY o.%s %s", s.Sort.By, s.Sort.Order)
	}
	return fmt.Sprintf(`SELECT o.%s, cb.name as created_by_name, ub.name as updated_by_name
		FROM %s AS o %s
    JOIN users cb on cb.id = o.created_by
		JOIN users ub on ub.id = o.updated_by
		WHERE TRUE %s AND o.deleted_at IS NULL
		%s
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", o."), c.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

func (r *orderRepo) Search(ctx context.Context, s *SearchOrdersOpts) ([]*OrderData, error) {
	orders := make([]*OrderData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&orders)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return orders, nil
}
func (r *orderRepo) Count(ctx context.Context, s *SearchOrdersOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	fmt.Println(sql)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("attendances.Count: %w", err)
	}

	return countResult, nil
}
