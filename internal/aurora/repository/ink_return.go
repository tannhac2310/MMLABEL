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

type SearchInkReturnOpts struct {
	ID                string
	InkID             string
	ProductionOrderID string
	Name              string
	Status            enum.InventoryCommonStatus
	Limit             int64
	Offset            int64
	Sort              *Sort
}

type InkReturnData struct {
	*model.InkReturn
}

// InkReturnRepo is a repository interface for inkReturn
type InkReturnRepo interface {
	Insert(ctx context.Context, e *model.InkReturn) error
	Update(ctx context.Context, e *model.InkReturn) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchInkReturnOpts) ([]*InkReturnData, error)
	Count(ctx context.Context, s *SearchInkReturnOpts) (*CountResult, error)
}

type inkReturnRepo struct {
}

func (i *inkReturnRepo) Insert(ctx context.Context, e *model.InkReturn) error {
	// insert to inkReturn
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}
	return nil
}

func (i *inkReturnRepo) Update(ctx context.Context, e *model.InkReturn) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (i *inkReturnRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE ink_return SET deleted_at = NOW() WHERE id = $1`
	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found any records to delete")
	}
	return nil
}

func (s *SearchInkReturnOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if s.ID != "" {
		conds += " AND b.id = $1"
		args = append(args, s.ID)
	}

	if s.Name != "" {
		args = append(args, "%"+s.Name+"%")
		conds += fmt.Sprintf(" AND( b.%[1]s ILIKE $%[3]d OR  b.%[2]s ILIKE $%[3]d)", model.InkReturnFieldName, model.InkReturnFieldCode, len(args))
	}

	if s.Status > 0 {
		args = append(args, s.Status)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.InkReturnFieldStatus, len(args))
	}

	if s.ID != "" {
		args = append(args, s.ID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.InkReturnFieldID, len(args))
	}

	b := &model.InkReturn{}
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

func (i *inkReturnRepo) Search(ctx context.Context, s *SearchInkReturnOpts) ([]*InkReturnData, error) {
	inkReturnData := make([]*InkReturnData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&inkReturnData)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return inkReturnData, nil
}

func (i *inkReturnRepo) Count(ctx context.Context, s *SearchInkReturnOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("chat.Count: %w", err)
	}

	return countResult, nil
}

// NewInkReturnRepo is a constructor for inkReturn repository
func NewInkReturnRepo() InkReturnRepo {
	return &inkReturnRepo{}
}
