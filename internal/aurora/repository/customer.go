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
	FindByID(ctx context.Context, id string) (*CustomerData, error)
	Search(ctx context.Context, s *SearchCustomerOpts) ([]*CustomerData, error)
	Count(ctx context.Context, s *SearchCustomerOpts) (*CountResult, error)
}

type sCustomerRepo struct {
}

func NewCustomerRepo() CustomerRepo {
	return &sCustomerRepo{}
}

func (r *sCustomerRepo) Insert(ctx context.Context, e *model.Customer) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sCustomerRepo) FindByID(ctx context.Context, id string) (*CustomerData, error) {
	e := &CustomerData{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("sCustomerRepo.cockroach.FindOne: %w", err)
	}

	return e, nil
}
func (r *sCustomerRepo) Update(ctx context.Context, e *model.Customer) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *sCustomerRepo) SoftDelete(ctx context.Context, id string) error {
	sql := "UPDATE customers SET deleted_at = NOW() WHERE id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("customers cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*sCustomerRepo not found any records to delete")
	}

	return nil
}

// SearchCustomerOpts all params is options
type SearchCustomerOpts struct {
	IDs    []string
	Name   string
	Code   string
	Phone  string
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchCustomerOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.CustomerFieldID)
	}
	// todo add more search options example:
	if s.Name != "" {
		args = append(args, "%"+s.Name+"%")
		conds += fmt.Sprintf(" AND (b.%[2]s ILIKE $%[1]d OR b.%[3]s ILIKE $%[1]d)",
			len(args), model.CustomerFieldName, model.CustomerFieldCode)
	}
	if s.Code != "" {
		args = append(args, s.Code)
		conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.CustomerFieldCode, len(args))
	}
	if s.Phone != "" {
		args = append(args, s.Phone)
		conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.CustomerFieldPhoneNumber, len(args))
	}

	b := &model.Customer{}
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

type CustomerData struct {
	*model.Customer
}

func (r *sCustomerRepo) Search(ctx context.Context, s *SearchCustomerOpts) ([]*CustomerData, error) {
	Customer := make([]*CustomerData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&Customer)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return Customer, nil
}

func (r *sCustomerRepo) Count(ctx context.Context, s *SearchCustomerOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sCustomerRepo.Count: %w", err)
	}

	return countResult, nil
}
