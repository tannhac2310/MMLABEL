package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type DepartmentRepo interface {
	Insert(ctx context.Context, e *model.Department) error
	Update(ctx context.Context, e *model.Department) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchDepartmentsOpts) ([]*DepartmentData, error)
	Count(ctx context.Context, s *SearchDepartmentsOpts) (*CountResult, error)
}

type departmentsRepo struct {
}

func NewDepartmentRepo() DepartmentRepo {
	return &departmentsRepo{}
}

func (r *departmentsRepo) Insert(ctx context.Context, e *model.Department) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *departmentsRepo) Update(ctx context.Context, e *model.Department) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *departmentsRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE departments
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

// SearchDepartmentsOpts all params is options
type SearchDepartmentsOpts struct {
	IDs    []string
	Name   string
	Code   string
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchDepartmentsOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.DepartmentFieldID)
	}
	if s.Name != "" {
		args = append(args, s.Name)
		conds += fmt.Sprintf(" AND (b.%[2]s ILIKE $%[1]d OR b.%[2]s ILIKE $%[1]d OR b.%[3]s ILIKE $%[1]d)",
			len(args), model.DepartmentFieldName, model.DepartmentFieldShortName, model.DepartmentFieldCode)
	}
	if s.Code != "" {
		args = append(args, s.Code)
		conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.DepartmentFieldCode, len(args))
	}

	b := &model.Department{}
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

type DepartmentData struct {
	*model.Department
}

func (r *departmentsRepo) Search(ctx context.Context, s *SearchDepartmentsOpts) ([]*DepartmentData, error) {
	message := make([]*DepartmentData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&message)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return message, nil
}

func (r *departmentsRepo) Count(ctx context.Context, s *SearchDepartmentsOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("chat.Count: %w", err)
	}

	return countResult, nil
}
