package repository

import (
	"context"
	"fmt"
	"strings"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type HistoryRepo interface {
	Insert(ctx context.Context, e *model.History) error
	Search(ctx context.Context, s *SearchHistoryOpts) ([]*HistoryData, error)
}

type historyRepo struct {
}

func NewHistoryRepo() HistoryRepo {
	return &historyRepo{}
}

func (r *historyRepo) Insert(ctx context.Context, e *model.History) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

type SearchHistoryOpts struct {
	IDs    []string
	Table  string
	Code   string
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchHistoryOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.HistoryFieldID)
	}

	if s.Code != "" {
		args = append(args, s.Table)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.HistoryFieldTable, len(args))
	}

	b := &model.History{}
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

type HistoryData struct {
	*model.History
}

func (r *historyRepo) Search(ctx context.Context, s *SearchHistoryOpts) ([]*HistoryData, error) {
	historys := make([]*HistoryData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&historys)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return historys, nil
}
