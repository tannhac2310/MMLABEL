package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type ZaloRepo interface {
	Insert(ctx context.Context, e *model.Zalo) error
	Update(ctx context.Context, e *model.Zalo) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchZalosOpts) ([]*ZaloData, error)
	Count(ctx context.Context, s *SearchZalosOpts) (*CountResult, error)
}

type zalosRepo struct {
}

func NewZaloRepo() ZaloRepo {
	return &zalosRepo{}
}

func (r *zalosRepo) Insert(ctx context.Context, e *model.Zalo) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *zalosRepo) Update(ctx context.Context, e *model.Zalo) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *zalosRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE zalo
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

// SearchZalosOpts all params is options
type SearchZalosOpts struct {
	IDs    []string
	OaID   string
	AppID  string
	Search string
	UserID string
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchZalosOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""
	if s.UserID != "" {
		joins = permissionCondition(enum.PermissionEntityOa, s.UserID, "b")
	}

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.ZaloFieldID)
	}

	if s.OaID != "" {
		args = append(args, s.OaID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.ZaloFieldOaID, len(args))
	}

	if s.AppID != "" {
		args = append(args, s.AppID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.ZaloFieldAppID, len(args))
	}
	if s.Search != "" {
		args = append(args, "%"+s.Search+"%")
		// oaId
		conds += fmt.Sprintf(" AND (b.%s ILIKE $%d", model.ZaloFieldOaID, len(args))
		// oaName
		conds += fmt.Sprintf(" OR b.%s ILIKE $%d )", model.ZaloFieldOaName, len(args))
	}

	b := &model.Zalo{}
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

type ZaloData struct {
	*model.Zalo
	CreatedByName string `db:"created_by_name"`
	UpdatedByName string `db:"updated_by_name"`
}

func (r *zalosRepo) Search(ctx context.Context, s *SearchZalosOpts) ([]*ZaloData, error) {
	zalo := make([]*ZaloData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&zalo)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return zalo, nil
}

func (r *zalosRepo) Count(ctx context.Context, s *SearchZalosOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("zalos.Count: %w", err)
	}

	return countResult, nil
}
