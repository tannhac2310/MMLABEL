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

type OptionRepo interface {
	Insert(ctx context.Context, e *model.Option) error
	Update(ctx context.Context, e *model.Option) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchOptionsOpts) ([]*OptionData, error)
	Count(ctx context.Context, s *SearchOptionsOpts) (*CountResult, error)
	FindByID(ctx context.Context, id string) (*OptionData, error)
}

type optionRepo struct {
}

func (r *optionRepo) Insert(ctx context.Context, e *model.Option) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *optionRepo) Update(ctx context.Context, e *model.Option) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *optionRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE options
		SET deleted_at = NOW()
		WHERE id = $1`

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found any records to delete")
	}

	return nil
}
func (i *optionRepo) FindByID(ctx context.Context, id string) (*OptionData, error) {
	optionData := &OptionData{}
	sql := `SELECT b.* FROM options AS b WHERE b.id = $1 AND b.deleted_at IS NULL`
	err := cockroach.Select(ctx, sql, id).ScanOne(optionData)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}
	return optionData, nil
}
// SearchOptionsOpts all params is options
type SearchOptionsOpts struct {
	IDs    []string
	Entity   string
	Name   string
	Code   string
	Status enum.CommonStatus
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchOptionsOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.OptionFieldID)
	}
	if s.Name != "" {
		args = append(args, "%"+s.Name+"%")
		conds += fmt.Sprintf(" AND (b.%[2]s ILIKE $%[1]d OR b.%[3]s ILIKE $%[1]d)",
			len(args), model.OptionFieldName, model.OptionFieldCode)
	}
	if s.Code != "" {
		args = append(args, s.Code)
		conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.OptionFieldCode, len(args))
	}
	if s.Entity != "" {
		args = append(args, s.Entity)
		conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.OptionFieldEntity, len(args))
	}

	b := &model.Option{}
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

type OptionData struct {
	*model.Option
}

func (r *optionRepo) Search(ctx context.Context, s *SearchOptionsOpts) ([]*OptionData, error) {
	options := make([]*OptionData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&options)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return options, nil
}

func (r *optionRepo) Count(ctx context.Context, s *SearchOptionsOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("chat.Count: %w", err)
	}

	return countResult, nil
}

func NewOptionRepo() OptionRepo {
	return &optionRepo{}
}
