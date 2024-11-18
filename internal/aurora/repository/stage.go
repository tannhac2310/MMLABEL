package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type StageRepo interface {
	Insert(ctx context.Context, e *model.Stage) error
	Update(ctx context.Context, e *model.Stage) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchStagesOpts) ([]*StageData, error)
	Count(ctx context.Context, s *SearchStagesOpts) (*CountResult, error)
}

type stagesRepo struct {
}

func NewStageRepo() StageRepo {
	return &stagesRepo{}
}

func (r *stagesRepo) Insert(ctx context.Context, e *model.Stage) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *stagesRepo) Update(ctx context.Context, e *model.Stage) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *stagesRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE stages
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

// SearchStagesOpts all params is options
type SearchStagesOpts struct {
	IDs    []string
	Name   string
	Code   string
	Codes  []string
	UserID string
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchStagesOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.StageFieldID)
	}

	if s.Name != "" {
		args = append(args, "%"+s.Name+"%")
		conds += fmt.Sprintf(" AND (b.%[1]s ILIKE $%[3]d OR b.%[2]s ILIKE $%[3]d)", model.StageFieldName, model.StageFieldShortName, len(args))
	}

	if s.Code != "" {
		args = append(args, s.Code)
		conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.StageFieldCode, len(args))
	}

	if len(s.Codes) > 0 {
		args = append(args, s.Codes)
		conds += fmt.Sprintf(" AND b.%s = ANY($%d)", model.StageFieldCode, len(args))
	}
	fmt.Println("============================================>>s.UserID", s.UserID)
	if s.UserID != "" {
		args = append(args, s.UserID)
		// join user_role, role_permission
		joins += fmt.Sprintf(` INNER JOIN user_role AS ur ON ur.user_id = $%d AND ur.deleted_at IS NULL
			INNER JOIN role_permissions AS rp ON rp.role_id = ur.role_id AND rp.entity_type = 'stage' and rp.entity_id = b.id`, len(args))
	}

	b := &model.Stage{}
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

	return fmt.Sprintf(`SELECT DISTINCT b.%s
		FROM %s AS b %s
		WHERE TRUE %s AND b.deleted_at IS NULL
		%s
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type StageData struct {
	*model.Stage
}

func (r *stagesRepo) Search(ctx context.Context, s *SearchStagesOpts) ([]*StageData, error) {
	message := make([]*StageData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&message)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return message, nil
}

func (r *stagesRepo) Count(ctx context.Context, s *SearchStagesOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("chat.Count: %w", err)
	}

	return countResult, nil
}
