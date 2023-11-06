package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type CustomerRepo interface {
	Insert(ctx context.Context, e *model.Customer) error
	Update(ctx context.Context, e *model.Customer) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchCustomersOpts) ([]*CustomerData, error)
	Count(ctx context.Context, s *SearchCustomersOpts) (*CountResult, error)
	SearchOne(ctx context.Context, s *SearchCustomersOpts) (*CustomerData, error)
}

type customersRepo struct {
}

func NewCustomerRepo() CustomerRepo {
	return &customersRepo{}
}

func (r *customersRepo) Insert(ctx context.Context, e *model.Customer) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *customersRepo) Update(ctx context.Context, e *model.Customer) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *customersRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE customers
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

// SearchCustomersOpts all params is options
type SearchCustomersOpts struct {
	IDs    []string
	Name   string
	Phone  string
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchCustomersOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.CustomerFieldID)
	}
	if s.Name != "" {
		args = append(args, s.Name)
		conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.CustomerFieldName, len(args))
	}
	if s.Phone != "" {
		args = append(args, s.Phone)
		conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.CustomerFieldPhoneNumber, len(args))
	}

	b := &model.Customer{}
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

type CustomerData struct {
	*model.Customer
}

func (r *customersRepo) Search(ctx context.Context, s *SearchCustomersOpts) ([]*CustomerData, error) {
	message := make([]*CustomerData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&message)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return message, nil
}
func (r *customersRepo) SearchOne(ctx context.Context, s *SearchCustomersOpts) (*CustomerData, error) {
	message := &CustomerData{}
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanOne(message)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select1: %w", err)
	}

	return message, nil
}

func (r *customersRepo) Count(ctx context.Context, s *SearchCustomersOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("chat.Count: %w", err)
	}

	return countResult, nil
}
