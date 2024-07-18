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

type SearchInkExportOpts struct {
	ID                  string
	Search              string
	ProductionOrderID   string
	ProductionOrderName string
	InkCode             string
	Status              enum.InventoryCommonStatus
	ExportDateFrom      time.Time
	ExportDateTo        time.Time
	Limit               int64
	Offset              int64
	Sort                *Sort
}

type InkExportData struct {
	*model.InkExport
	CreatedByName string `db:"created_by_name"`
	UpdatedByName string `db:"updated_by_name"`
}

// InkExportRepo is a repository interface for inkExport
type InkExportRepo interface {
	Insert(ctx context.Context, e *model.InkExport) error
	Update(ctx context.Context, e *model.InkExport) error
	SoftDelete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.InkExport, error)
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
	sql := `UPDATE ink_export SET deleted_at = NOW() WHERE id = $1`
	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found any records to delete")
	}
	return nil
}

func (i *inkExportRepo) FindByID(ctx context.Context, id string) (*model.InkExport, error) {
	e := &model.InkExport{}
	if err := cockroach.FindOne(ctx, e, "id = $1", id); err != nil {
		return nil, err
	}

	return e, nil
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

// buildSearchInkExportQuery is a helper function to build query for search inkExports
func (i *SearchInkExportOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := " JOIN users AS cu ON cu.id = b.created_by " +
		"JOIN users AS uu ON uu.id = b.updated_by "

	if i.ID != "" {
		conds += " AND b.id = $1"
		args = append(args, i.ID)
	}

	if i.ProductionOrderName != "" {
		joins += fmt.Sprintf(" LEFT JOIN production_orders AS po ON po.id = b.%s", model.InkExportFieldProductionOrderID)
		args = append(args, "%"+i.ProductionOrderName+"%")
		conds += fmt.Sprintf(" AND po.%s ILIKE $%d ", model.ProductionOrderFieldName, len(args))
	}
	if i.InkCode != "" {
		joins += fmt.Sprintf(" LEFT JOIN ink_export_detail AS id ON id.ink_export_id = b.%s", model.InkExportFieldID)
		joins += fmt.Sprintf(" LEFT JOIN ink ON ink.id = id.%s", model.InkExportDetailFieldInkID)
		args = append(args, "%"+i.InkCode+"%")
		conds += fmt.Sprintf(" AND ink.%s ILIKE $%d ", model.InkFieldCode, len(args))
	}

	if i.ProductionOrderID != "" {
		joins += fmt.Sprintf(" LEFT JOIN ink AS ik ON ik.id = b.%s LEFT JOIN ", model.InkExportDetailFieldInkID)
		args = append(args, i.ProductionOrderID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.InkExportFieldProductionOrderID, len(args))
	}

	if !i.ExportDateFrom.IsZero() {
		args = append(args, i.ExportDateFrom)
		conds += fmt.Sprintf(" AND b.%s >= $%d", model.InkExportFieldExportDate, len(args))
	}
	if !i.ExportDateTo.IsZero() {
		args = append(args, i.ExportDateTo)
		conds += fmt.Sprintf(" AND b.%s <= $%d", model.InkExportFieldExportDate, len(args))
	}
	if i.Search != "" {
		joins += fmt.Sprintf(" LEFT JOIN production_orders AS po ON po.id = b.%s", model.InkExportFieldProductionOrderID)
		args = append(args, "%"+i.Search+"%")
		conds += fmt.Sprintf(" AND (b.%[2]s ILIKE $%[1]d OR b.%[3]s ILIKE $%[1]d OR b.%[4]s ILIKE $%[1]d OR b.%[5]s ILIKE $%[1]d OR po.name ILIKE $%[1]d)",
			len(args), model.InkExportFieldCode, model.InkExportFieldName, model.InkExportFieldDescription, model.InkExportFieldProductionOrderID,
		)
	}

	if i.Status > 0 {
		args = append(args, i.Status)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.InkExportFieldStatus, len(args))
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
	return fmt.Sprintf(`SELECT b.%s, cu.name as created_by_name, uu.name as updated_by_name
		FROM %s AS b %s
		WHERE TRUE %s AND b.deleted_at IS NULL
		%s
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", b."), b.TableName(), joins, conds, order, i.Limit, i.Offset), args
}

// NewInkExportRepo is a constructor for inkExport repository
func NewInkExportRepo() InkExportRepo {
	return &inkExportRepo{}
}
