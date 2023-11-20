package repository

import (
	"context"
	"fmt"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"strings"
	"time"
)

type SearchInkExportOpts struct {
	ID                string
	Name              string
	ProductionOrderID string
	Code              string
	Limit             int64
	Offset            int64
	Sort              *Sort
}

type InkExportData struct {
	*model.InkExport
}

// InkExportRepo is a repository interface for inkExport
type InkExportRepo interface {
	Insert(ctx context.Context, e *model.InkExport) error
	Update(ctx context.Context, e *model.InkExport) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchInkExportOpts) ([]*InkExportData, error)
	Count(ctx context.Context, s *SearchInkExportOpts) (*CountResult, error)
}

type inkExportRepo struct {
}

func (i *inkExportRepo) Insert(ctx context.Context, e *model.InkExport) error {
	// insert to inkExport
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}
	return nil
}

func (i *inkExportRepo) Update(ctx context.Context, e *model.InkExport) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (i *inkExportRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE inkExport SET deleted_at = NOW() WHERE id = $1`
	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found any records to delete")
	}
	return nil
}

// buildSearchInkExportQuery is a helper function to build query for search inkExports
func (i *SearchInkExportOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if i.ID != "" {
		conds += " AND b.id = $1"
		args = append(args, i.ID)
	}
	if i.ProductionOrderID != "" {
		conds += " AND b.production_order_id = $1"
		args = append(args, i.ProductionOrderID)
	}
	if i.Name != "" {
		conds += " AND b.name ILIKE $1"
		args = append(args, "%"+i.Name+"%")
	}

	if i.Code != "" {
		conds += " AND b.code ILIKE $1"
		args = append(args, "%"+i.Code+"%")
	}

	b := &model.InkExport{}
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
		WHERE TRUE %s AND b.deleted_at IS NULL
		%s
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", b."), b.TableName(), joins, conds, order, i.Limit, i.Offset), args
}

func (i *inkExportRepo) Search(ctx context.Context, s *SearchInkExportOpts) ([]*InkExportData, error) {
	inkExportData := make([]*InkExportData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&inkExportData)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return inkExportData, nil
}

func (i *inkExportRepo) Count(ctx context.Context, s *SearchInkExportOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("chat.Count: %w", err)
	}

	return countResult, nil
}

// NewInkExportRepo is a constructor for inkExport repository
func NewInkExportRepo() InkExportRepo {
	return &inkExportRepo{}
}
