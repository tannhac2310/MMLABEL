package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

type RoleRepo interface {
	Insert(ctx context.Context, e *model.Role) error
	Update(ctx context.Context, e *model.Role) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.Role, error)
	HighestRole(ctx context.Context, id []string) (*model.Role, error)
	Search(ctx context.Context, s *SearchRoleOpts) ([]*model.Role, error)
}

type roleRepo struct {
}

func NewRoleRepo() RoleRepo {
	return &roleRepo{}
}

func (r *roleRepo) Insert(ctx context.Context, e *model.Role) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("cockroach.Create: %w", err)
	}

	return nil
}
func (r *roleRepo) Delete(ctx context.Context, id string) error {
	sql := `DELETE FROM roles
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
func (r *roleRepo) FindByID(ctx context.Context, id string) (*model.Role, error) {
	e := &model.Role{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("cockroach.FindOne: %w", err)
	}

	return e, nil
}

func (r *roleRepo) Update(ctx context.Context, e *model.Role) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *roleRepo) HighestRole(ctx context.Context, id []string) (*model.Role, error) {
	e := &model.Role{}
	fields, values := e.FieldMap()
	sql := fmt.Sprintf(`SELECT %s
		FROM %s
		WHERE id = ANY($1)
		ORDER BY priority ASC`,
		strings.Join(fields, ","), e.TableName())

	err := cockroach.QueryRow(ctx, sql, id).Scan(values...)
	if err != nil {
		return nil, fmt.Errorf("cockroach.QueryRow: %w", err)
	}

	return e, nil
}

type SearchRoleOpts struct {
	IDs    []string
	Name   string
	Limit  int64
	Offset int64
}

func (r *roleRepo) Search(ctx context.Context, s *SearchRoleOpts) ([]*model.Role, error) {
	ponds := make([]*model.Role, 0)
	sql, args := s.buildQuery()

	err := cockroach.Select(ctx, sql, args...).ScanAll(&ponds)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return ponds, nil
}

func (s *SearchRoleOpts) buildQuery() (string, []interface{}) {
	args := []interface{}{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND r.%s = ANY($1)", model.RoleFieldID)
	}

	if s.Name != "" {
		args = append(args, "%"+s.Name+"%")
		conds += fmt.Sprintf(" AND r.%s ILIKE $%d", model.RoleFieldName, len(args))
	}

	e := &model.Role{}
	fields, _ := e.FieldMap()

	return fmt.Sprintf(`SELECT r.%s
		FROM %s AS r %s
		WHERE TRUE %s AND r.deleted_at IS NULL
		ORDER BY r.id DESC
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", r."), e.TableName(), joins, conds, s.Limit, s.Offset), args
}
