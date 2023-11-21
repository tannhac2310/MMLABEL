package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type SearchInkExportDetailOpts struct {
	ProductionOrderID string
	InkExportID       string
	InkID             string
	Limit             int64
	Offset            int64
	Sort              *Sort
}

type InkExportDetailData struct {
	*model.InkExportDetail
}

// InkExportDetailRepo is a repository interface for inkExportDetail
type InkExportDetailRepo interface {
	Insert(ctx context.Context, e *model.InkExportDetail) error
	Update(ctx context.Context, e *model.InkExportDetail) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchInkExportDetailOpts) ([]*InkExportDetailData, error)
	Count(ctx context.Context, s *SearchInkExportDetailOpts) (*CountResult, error)
}

type inkExportDetailRepo struct {
}

func (i *inkExportDetailRepo) Insert(ctx context.Context, e *model.InkExportDetail) error {
	// insert to inkExportDetail
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}
	return nil
}

func (i *inkExportDetailRepo) Update(ctx context.Context, e *model.InkExportDetail) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (i *inkExportDetailRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE inK_export_detail SET deleted_at = NOW() WHERE id = $1`
	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found any records to delete")
	}
	return nil
}

// buildSearchInkExportDetailQuery is a helper function to build query for search inkExportDetails
func (i *SearchInkExportDetailOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if i.InkExportID != "" {
		conds += " AND b.ink_export_id = $1"
		args = append(args, i.InkExportID)
	}
	if i.InkID != "" {
		args = append(args, i.InkID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.InkExportDetailFieldInkID, len(args))
	}
	if i.ProductionOrderID != "" {
		args = append(args, i.ProductionOrderID)
		conds += fmt.Sprintf(" AND ie.%s = $%d", model.InkExportFieldProductionOrderID, len(args))
	}

	b := &model.InkExportDetail{}
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
		JOIN ink_export AS ie ON ie.id = b.ink_export_id
		WHERE TRUE %s AND b.deleted_at IS NULL
		%s
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", b."), b.TableName(), joins, conds, order, i.Limit, i.Offset), args
}

func (i *inkExportDetailRepo) Search(ctx context.Context, s *SearchInkExportDetailOpts) ([]*InkExportDetailData, error) {
	inkExportDetailData := make([]*InkExportDetailData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&inkExportDetailData)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return inkExportDetailData, nil
}

func (i *inkExportDetailRepo) Count(ctx context.Context, s *SearchInkExportDetailOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("chat.Count: %w", err)
	}

	return countResult, nil
}

// NewInkExportDetailRepo is a constructor for inkExportDetail repository
func NewInkExportDetailRepo() InkExportDetailRepo {
	return &inkExportDetailRepo{}
}
