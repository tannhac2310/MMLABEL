package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type SearchInkReturnDetailOpts struct {
	InkReturnID string
	InkID       string
	InkCode     string
	InkExportID string
	Limit       int64
	Offset      int64
	Sort        *Sort
}

type InkReturnDetailData struct {
	*model.InkReturnDetail
}

// InkReturnDetailRepo is a repository interface for inkReturnDetail
type InkReturnDetailRepo interface {
	Insert(ctx context.Context, e *model.InkReturnDetail) error
	Update(ctx context.Context, e *model.InkReturnDetail) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchInkReturnDetailOpts) ([]*InkReturnDetailData, error)
	Count(ctx context.Context, s *SearchInkReturnDetailOpts) (*CountResult, error)
}

type inkReturnDetailRepo struct {
}

func (i *inkReturnDetailRepo) Insert(ctx context.Context, e *model.InkReturnDetail) error {
	// insert to inkReturnDetail
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}
	return nil
}

func (i *inkReturnDetailRepo) Update(ctx context.Context, e *model.InkReturnDetail) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (i *inkReturnDetailRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE ink_return_detail SET deleted_at = NOW() WHERE id = $1`
	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found any records to delete")
	}
	return nil
}

func (i *inkReturnDetailRepo) Search(ctx context.Context, s *SearchInkReturnDetailOpts) ([]*InkReturnDetailData, error) {
	inkReturnDetailData := make([]*InkReturnDetailData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&inkReturnDetailData)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return inkReturnDetailData, nil
}

func (i *inkReturnDetailRepo) Count(ctx context.Context, s *SearchInkReturnDetailOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("chat.Count: %w", err)
	}

	return countResult, nil
}

// buildSearchInkReturnDetailQuery is a helper function to build query for search inkReturnDetails
func (i *SearchInkReturnDetailOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if i.InkReturnID != "" {
		args = append(args, i.InkReturnID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.InkReturnDetailFieldInkReturnID, len(args))
	}
	if i.InkID != "" {
		args = append(args, i.InkID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.InkReturnDetailFieldInkID, len(args))
	}
	if i.InkCode != "" {
		args = append(args, i.InkCode)
		conds += fmt.Sprintf(" AND i.%s = $%d", model.InkFieldCode, len(args))
	}
	if i.InkExportID != "" {
		args = append(args, i.InkExportID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.InkReturnDetailFieldInkExportID, len(args))
	}

	b := &model.InkReturnDetail{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf(`SELECT count(*) as cnt
		FROM %s AS b %s
		WHERE TRUE %s AND b.deleted_at IS NULL`, b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.id DESC "
	if i.Sort != nil {
		order = fmt.Sprintf(" ORDER BY b.%s %s", i.Sort.By, i.Sort.Order)
	}
	return fmt.Sprintf(`SELECT b.%s
		FROM %s AS b %s
		JOIN ink AS i ON i.id = b.ink_id
		WHERE TRUE %s AND b.deleted_at IS NULL
		%s
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", b."), b.TableName(), joins, conds, order, i.Limit, i.Offset), args
}

// NewInkReturnDetailRepo is a constructor for inkReturnDetail repository
func NewInkReturnDetailRepo() InkReturnDetailRepo {
	return &inkReturnDetailRepo{}
}
