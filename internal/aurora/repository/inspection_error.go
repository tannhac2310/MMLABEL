package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type InspectionErrorRepo interface {
	Insert(ctx context.Context, e *model.InspectionError) error
	Update(ctx context.Context, e *model.InspectionError) error
	SoftDeleteByFormID(ctx context.Context, formID string) error
	FindByID(ctx context.Context, id string) (*InspectionErrorData, error)
	Search(ctx context.Context, s *SearchInspectionErrorOpts) ([]*InspectionErrorData, error)
	Count(ctx context.Context, s *SearchInspectionErrorOpts) (*CountResult, error)
}

type sInspectionErrorRepo struct {
}

func NewInspectionErrorRepo() InspectionErrorRepo {
	return &sInspectionErrorRepo{}
}

func (r *sInspectionErrorRepo) Insert(ctx context.Context, e *model.InspectionError) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sInspectionErrorRepo) FindByID(ctx context.Context, id string) (*InspectionErrorData, error) {
	e := &InspectionErrorData{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("sInspectionErrorRepo.cockroach.FindOne: %w", err)
	}

	return e, nil
}
func (r *sInspectionErrorRepo) Update(ctx context.Context, e *model.InspectionError) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *sInspectionErrorRepo) SoftDeleteByFormID(ctx context.Context, formID string) error {
	sql := "UPDATE inspection_errors SET deleted_at = NOW() WHERE inspection_form_id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, formID)
	if err != nil {
		return fmt.Errorf("inspection_errors cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*sInspectionErrorRepo not found any records to delete")
	}

	return nil
}

// SearchInspectionErrorOpts all params is options
type SearchInspectionErrorOpts struct {
	IDs               []string
	InspectionFormIDs []string
	Limit             int64
	Offset            int64
	Sort              *Sort
}

func (s *SearchInspectionErrorOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.InspectionErrorFieldID)
	}
	if len(s.InspectionFormIDs) > 0 {
		args = append(args, s.InspectionFormIDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($%d)", model.InspectionErrorFieldInspectionFormID, len(args))
	}

	b := &model.InspectionError{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf("SELECT count(*) as cnt FROM %s AS b %s WHERE TRUE %s AND b.deleted_at IS NULL", b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.id DESC "
	if s.Sort != nil {
		order = fmt.Sprintf(" ORDER BY b.%s %s", s.Sort.By, s.Sort.Order)
	}
	return fmt.Sprintf("SELECT b.%s FROM %s AS b %s WHERE TRUE %s AND b.deleted_at IS NULL %s LIMIT %d OFFSET %d", strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type InspectionErrorData struct {
	*model.InspectionError
}

func (r *sInspectionErrorRepo) Search(ctx context.Context, s *SearchInspectionErrorOpts) ([]*InspectionErrorData, error) {
	InspectionError := make([]*InspectionErrorData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&InspectionError)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return InspectionError, nil
}

func (r *sInspectionErrorRepo) Count(ctx context.Context, s *SearchInspectionErrorOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sInspectionErrorRepo.Count: %w", err)
	}

	return countResult, nil
}
