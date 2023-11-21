package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type SearchInkImportDetailOpts struct {
	InkImportID string
	InkID       string
	InkCode     string
	ID          string
	Limit       int64
	Offset      int64
	Sort        *Sort
}

type InkImportDetailData struct {
	*model.InkImportDetail
}

// InkImportDetailRepo is a repository interface for inkImportDetail
type InkImportDetailRepo interface {
	Insert(ctx context.Context, e *model.InkImportDetail) error
	Update(ctx context.Context, e *model.InkImportDetail) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchInkImportDetailOpts) ([]*InkImportDetailData, error)
	Count(ctx context.Context, s *SearchInkImportDetailOpts) (*CountResult, error)
}

type inkImportDetailRepo struct {
}

func (i *inkImportDetailRepo) Insert(ctx context.Context, e *model.InkImportDetail) error {
	// insert to inkImportDetail
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}
	return nil
}

func (i *inkImportDetailRepo) Update(ctx context.Context, e *model.InkImportDetail) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (i *inkImportDetailRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE ink_import_detail SET deleted_at = NOW() WHERE id = $1`
	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found any records to delete")
	}
	return nil
}

// buildSearchInkImportDetailQuery is a helper function to build query for search inkImportDetails
func (i *SearchInkImportDetailOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if i.InkImportID != "" {
		conds += " AND b.ink_import_id = $1"
		args = append(args, i.InkImportID)
	}

	if i.ID != "" {
		args = append(args, i.ID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.InkImportDetailFieldID, len(args))
	}

	if i.InkID != "" {
		args = append(args, i.InkID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.InkImportDetailFieldID, len(args)) // write ink.ID = ink_import.ID
	}

	if i.InkCode != "" {
		args = append(args, i.InkCode)
		conds += fmt.Sprintf(" AND i.%s = $%d", model.InkFieldCode, len(args))
	}

	b := &model.InkImportDetail{}
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
	// JOIN ink AS i ON i.id = b.id when importing, I write ink.ID = ink_import.ID
	return fmt.Sprintf(`SELECT b.%s
		FROM %s AS b %s
		JOIN ink AS i ON i.id = b.id 
		WHERE TRUE %s AND b.deleted_at IS NULL
		%s
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", b."), b.TableName(), joins, conds, order, i.Limit, i.Offset), args
}

func (i *inkImportDetailRepo) Search(ctx context.Context, s *SearchInkImportDetailOpts) ([]*InkImportDetailData, error) {
	inkImportDetailData := make([]*InkImportDetailData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&inkImportDetailData)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return inkImportDetailData, nil
}

func (i *inkImportDetailRepo) Count(ctx context.Context, s *SearchInkImportDetailOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("chat.Count: %w", err)
	}

	return countResult, nil
}

// NewInkImportDetailRepo is a constructor for inkImportDetail repository
func NewInkImportDetailRepo() InkImportDetailRepo {
	return &inkImportDetailRepo{}
}
