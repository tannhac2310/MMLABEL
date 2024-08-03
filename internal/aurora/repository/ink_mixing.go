package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type InkMixingRepo interface {
	Insert(ctx context.Context, e *model.InkMixing) error
	Update(ctx context.Context, e *model.InkMixing) error
	SoftDelete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.InkMixing, error)
	Search(ctx context.Context, s *SearchInkMixingOpts) ([]*InkMixingData, error)
	Count(ctx context.Context, s *SearchInkMixingOpts) (*CountResult, error)
}

type sInkMixingRepo struct {
}

func NewInkMixingRepo() InkMixingRepo {
	return &sInkMixingRepo{}
}

func (r *sInkMixingRepo) Insert(ctx context.Context, e *model.InkMixing) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sInkMixingRepo) FindByID(ctx context.Context, id string) (*model.InkMixing, error) {
	e := &model.InkMixing{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("sInkMixingRepo.cockroach.FindOne: %w", err)
	}
	return e, nil
}
func (r *sInkMixingRepo) Update(ctx context.Context, e *model.InkMixing) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *sInkMixingRepo) SoftDelete(ctx context.Context, id string) error {
	sql := "UPDATE ink_mixing SET deleted_at = NOW() WHERE id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("ink_mixing cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*sInkMixingRepo not found any records to delete")
	}

	return nil
}

// SearchInkMixingOpts all params is options
type SearchInkMixingOpts struct {
	IDs        []string
	MixingDate string
	Limit      int64
	Offset     int64
	Sort       *Sort
}

func (s *SearchInkMixingOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := " LEFT JOIN users AS u ON u.id = b.created_by LEFT JOIN users AS u2 ON u2.id = b.updated_by "

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.InkMixingFieldID)
	}

	//if s.MixingDate != "" {
	//	args = append(args, s.MixingDate)
	//	conds += fmt.Sprintf(" AND b.%s = $%d", model.InkMixingFieldMixingDate, len(args))
	//}

	b := &model.InkMixing{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf("SELECT count(*) as cnt FROM %s AS b %s WHERE TRUE %s AND b.deleted_at IS NULL", b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.id DESC "
	if s.Sort != nil {
		order = fmt.Sprintf(" ORDER BY b.%s %s", s.Sort.By, s.Sort.Order)
	}
	return fmt.Sprintf(`
SELECT b.%s, u.name AS created_by_name, u2.name AS updated_by_name
FROM %s AS b %s WHERE TRUE %s AND b.deleted_at IS NULL %s LIMIT %d OFFSET %d
`, strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type InkMixingData struct {
	*model.InkMixing
	CreatedByName string `json:"created_by_name"`
	UpdatedByName string `json:"updated_by_name"`
}

func (r *sInkMixingRepo) Search(ctx context.Context, s *SearchInkMixingOpts) ([]*InkMixingData, error) {
	InkMixing := make([]*InkMixingData, 0)
	sql, args := s.buildQuery(false)
	fmt.Println("=============>>", sql, args)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&InkMixing)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return InkMixing, nil
}

func (r *sInkMixingRepo) Count(ctx context.Context, s *SearchInkMixingOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sInkMixingRepo.Count: %w", err)
	}

	return countResult, nil
}
