package repository

import (
	"context"
	"fmt"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
	"strings"
)

type RolePermissionRepo interface {
	Insert(ctx context.Context, e *model.RolePermission) error
	DeleteByRoleID(ctx context.Context, roleId string) error
	Search(ctx context.Context, s *SearchRolePermissionOpts) ([]*RolePermissionData, error)
}

type sRolePermissionRepo struct {
}

func NewRolePermissionRepo() RolePermissionRepo {
	return &sRolePermissionRepo{}
}

func (r *sRolePermissionRepo) Insert(ctx context.Context, e *model.RolePermission) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sRolePermissionRepo) DeleteByRoleID(ctx context.Context, roleId string) error {
	sql := "DELETE FROM role_permissions WHERE role_id = $1;"

	_, err := cockroach.Exec(ctx, sql, roleId)
	if err != nil {
		return fmt.Errorf("role_permissions cockroach.Exec: %w", err)
	}

	return nil
}

// SearchRolePermissionOpts all params is options
type SearchRolePermissionOpts struct {
	RoleID string
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchRolePermissionOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if s.RoleID != "" {
		args = append(args, s.RoleID)
		conds += fmt.Sprintf(" AND b.role_id = $%d ", len(args))
	}

	b := &model.RolePermission{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf("SELECT count(*) as cnt FROM %s AS b %s WHERE TRUE %s ", b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.id DESC "
	if s.Sort != nil {
		order = fmt.Sprintf(" ORDER BY b.%s %s", s.Sort.By, s.Sort.Order)
	}
	return fmt.Sprintf("SELECT b.%s FROM %s AS b %s WHERE TRUE %s %s LIMIT %d OFFSET %d", strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type RolePermissionData struct {
	*model.RolePermission
}

func (r *sRolePermissionRepo) Search(ctx context.Context, s *SearchRolePermissionOpts) ([]*RolePermissionData, error) {
	RolePermission := make([]*RolePermissionData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&RolePermission)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return RolePermission, nil
}
