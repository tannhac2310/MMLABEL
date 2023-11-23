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

type SearchInkOpts struct {
	ID           string
	Name         string
	Code         string
	Manufacturer string
	Expiration   string
	Status       enum.CommonStatus
	Limit        int64
	Offset       int64
	Sort         *Sort
}

type InkData struct {
	*model.Ink
}

// InkRepo is a repository interface for ink
type InkRepo interface {
	Insert(ctx context.Context, e *model.Ink) error
	Update(ctx context.Context, e *model.Ink) error
	SoftDelete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*InkData, error)
	Search(ctx context.Context, s *SearchInkOpts) ([]*InkData, error)
	Count(ctx context.Context, s *SearchInkOpts) (*CountResult, error)
}

type inkRepo struct {
}

func (i *inkRepo) Insert(ctx context.Context, e *model.Ink) error {
	// insert to ink
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}
	return nil
}

func (i *inkRepo) Update(ctx context.Context, e *model.Ink) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}
func (i *inkRepo) FindByID(ctx context.Context, id string) (*InkData, error) {
	inkData := &InkData{}
	sql := `SELECT b.* FROM ink AS b WHERE b.id = $1 AND b.deleted_at IS NULL`
	err := cockroach.Select(ctx, sql, id).ScanOne(inkData)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}
	return inkData, nil
}

func (i *inkRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE ink SET deleted_at = NOW() WHERE id = $1`
	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found any records to delete")
	}
	return nil
}

// buildSearchInkQuery is a helper function to build query for search inks
func (i *SearchInkOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if i.ID != "" {
		conds += " AND b.id = $1"
		args = append(args, i.ID)
	}

	if i.Name != "" {
		args = append(args, "%"+i.Name+"%")
		conds += fmt.Sprintf(" AND( b.%[1]s ILIKE $%[3]d OR  b.%[2]s ILIKE $%[3]d)", model.InkFieldName, model.InkFieldCode, len(args))
	}

	if i.Code != "" {
		args = append(args, i.Code)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.InkFieldCode, len(args))
	}

	if i.Manufacturer != "" {
		args = append(args, i.Manufacturer)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.InkFieldManufacturer, len(args))
	}

	if i.Expiration != "" {
		args = append(args, i.Expiration)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.InkFieldExpirationDate, len(args))
	}

	if i.Status > 0 {
		args = append(args, i.Status)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.InkFieldStatus, len(args))
	}

	b := &model.Ink{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf(`SELECT count(*) as cnt
		FROM %s AS b %s
		WHERE TRUE %s AND b.deleted_at IS NULL`, b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.created_at DESC "
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

func (i *inkRepo) Search(ctx context.Context, s *SearchInkOpts) ([]*InkData, error) {
	inkData := make([]*InkData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&inkData)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return inkData, nil
}

func (i *inkRepo) Count(ctx context.Context, s *SearchInkOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("chat.Count: %w", err)
	}

	return countResult, nil
}

// NewInkRepo is a constructor for ink repository
func NewInkRepo() InkRepo {
	return &inkRepo{}
}
