package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type InkMixingDetailRepo interface {
	Insert(ctx context.Context, e *model.InkMixingDetail) error
	Update(ctx context.Context, e *model.InkMixingDetail) error
	SoftDelete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*InkMixingDetailData, error)
	Search(ctx context.Context, s *SearchInkMixingDetailOpts) ([]*InkMixingDetailData, error)
	Count(ctx context.Context, s *SearchInkMixingDetailOpts) (*CountResult, error)
}

type sInkMixingDetailRepo struct {
}

func NewInkMixingDetailRepo() InkMixingDetailRepo {
	return &sInkMixingDetailRepo{}
}

func (r *sInkMixingDetailRepo) Insert(ctx context.Context, e *model.InkMixingDetail) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sInkMixingDetailRepo) FindByID(ctx context.Context, id string) (*InkMixingDetailData, error) {
	e := &InkMixingDetailData{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("sInkMixingDetailRepo.cockroach.FindOne: %w", err)
	}

	return e, nil
}
func (r *sInkMixingDetailRepo) Update(ctx context.Context, e *model.InkMixingDetail) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *sInkMixingDetailRepo) SoftDelete(ctx context.Context, id string) error {
	sql := "UPDATE ink_mixing_detail SET deleted_at = NOW() WHERE id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("ink_mixing_detail cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*sInkMixingDetailRepo not found any records to delete")
	}

	return nil
}

// SearchInkMixingDetailOpts all params is options
type SearchInkMixingDetailOpts struct {
	IDs          []string
	InkMixingID  string
	InkMixingIDs []string
	Limit        int64
	Offset       int64
	Sort         *Sort
}

func (s *SearchInkMixingDetailOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.InkMixingDetailFieldID)
	}
	if s.InkMixingID != "" {
		args = append(args, s.InkMixingID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.InkMixingDetailFieldInkMixingID, len(args))
	}

	if len(s.InkMixingIDs) > 0 {
		args = append(args, s.InkMixingIDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($%d)", model.InkMixingDetailFieldInkMixingID, len(args))
	}

	b := &model.InkMixingDetail{}
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

type InkMixingDetailData struct {
	*model.InkMixingDetail
}

func (r *sInkMixingDetailRepo) Search(ctx context.Context, s *SearchInkMixingDetailOpts) ([]*InkMixingDetailData, error) {
	InkMixingDetail := make([]*InkMixingDetailData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&InkMixingDetail)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return InkMixingDetail, nil
}

func (r *sInkMixingDetailRepo) Count(ctx context.Context, s *SearchInkMixingDetailOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sInkMixingDetailRepo.Count: %w", err)
	}

	return countResult, nil
}
