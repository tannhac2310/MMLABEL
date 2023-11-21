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

type SearchInkImportOpts struct {
	ID     string
	Name   string
	Status enum.InventoryCommonStatus
	Limit  int64
	Offset int64
	Sort   *Sort
}

type InkImportData struct {
	*model.InkImport
}

// InkImportRepo is a repository interface for inkImport
type InkImportRepo interface {
	Insert(ctx context.Context, e *model.InkImport) error
	Update(ctx context.Context, e *model.InkImport) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchInkImportOpts) ([]*InkImportData, error)
	Count(ctx context.Context, s *SearchInkImportOpts) (*CountResult, error)
}

type inkImportRepo struct {
}

func (i *inkImportRepo) Insert(ctx context.Context, e *model.InkImport) error {
	// insert to inkImport
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}
	return nil
}

func (i *inkImportRepo) Update(ctx context.Context, e *model.InkImport) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (i *inkImportRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE ink_import SET deleted_at = NOW() WHERE id = $1`
	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found any records to delete")
	}
	return nil
}

func (s *SearchInkImportOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if s.ID != "" {
		conds += " AND b.id = $1"
		args = append(args, s.ID)
	}

	if s.Name != "" {
		args = append(args, "%"+s.Name+"%")
		conds += fmt.Sprintf(" AND( b.%[1]s ILIKE $%[3]d OR  b.%[2]s ILIKE $%[3]d)", model.InkImportFieldName, model.InkImportFieldCode, len(args))
	}

	if s.Status > 0 {
		args = append(args, s.Status)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.InkImportFieldStatus, len(args))
	}

	b := &model.InkImport{}
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

func (i *inkImportRepo) Search(ctx context.Context, s *SearchInkImportOpts) ([]*InkImportData, error) {
	inkImportData := make([]*InkImportData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&inkImportData)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return inkImportData, nil
}

func (i *inkImportRepo) Count(ctx context.Context, s *SearchInkImportOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("chat.Count: %w", err)
	}

	return countResult, nil
}

// NewInkImportRepo is a constructor for inkImport repository
func NewInkImportRepo() InkImportRepo {
	return &inkImportRepo{}
}
