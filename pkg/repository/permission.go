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

type PermissionRepo interface {
	Insert(ctx context.Context, e *model.Permission) error
	Update(ctx context.Context, e *model.Permission) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchPermissionsOpts) ([]*PermissionData, error)
	Count(ctx context.Context, s *SearchPermissionsOpts) (*CountResult, error)
}

type permissionsRepo struct {
}

func NewPermissionRepo() PermissionRepo {
	return &permissionsRepo{}
}

func (r *permissionsRepo) Insert(ctx context.Context, e *model.Permission) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *permissionsRepo) Update(ctx context.Context, e *model.Permission) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *permissionsRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE permission
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

// SearchPermissionsOpts all params is options
type SearchPermissionsOpts struct {
	IDs       []string
	UserID    string
	Entity    enum.PermissionEntity
	ElementID string
	Limit     int64
	Offset    int64
	Sort      *Sort
}

func (s *SearchPermissionsOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.PermissionFieldID)
	}

	if s.UserID != "" {
		args = append(args, s.UserID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.PermissionFieldUserID, len(args))
	}
	if s.Entity > 0 {
		args = append(args, s.Entity)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.PermissionFieldEntity, len(args))
	}
	if s.ElementID != "" {
		args = append(args, s.ElementID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.PermissionFieldElementID, len(args))
	}

	b := &model.Permission{}
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
	return fmt.Sprintf(`SELECT b.%s, 
		cb.name as created_by_name, ub.name as updated_by_name 
		FROM %s AS b %s
		JOIN users cb on cb.id = b.created_by
		JOIN users ub on ub.id = b.updated_by
		WHERE TRUE %s AND b.deleted_at IS NULL
		%s
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type PermissionData struct {
	*model.Permission
	CreatedByName string `db:"created_by_name"`
	UpdatedByName string `db:"updated_by_name"`
}

func (r *permissionsRepo) Search(ctx context.Context, s *SearchPermissionsOpts) ([]*PermissionData, error) {
	permission := make([]*PermissionData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&permission)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return permission, nil
}

func (r *permissionsRepo) Count(ctx context.Context, s *SearchPermissionsOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("permissions.Count: %w", err)
	}

	return countResult, nil
}
