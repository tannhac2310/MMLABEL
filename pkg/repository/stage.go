package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

type StageRepo interface {
	Insert(ctx context.Context, e *model.Stage) error
	Update(ctx context.Context, e *model.Stage) error
	FindByID(ctx context.Context, id string) (*model.Stage, error)
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchStagesOpts) ([]*model.Stage, error)
	Count(ctx context.Context, s *SearchStagesOpts) (*CountResult, error)
}

type stageRepo struct {
}

func NewStageRepo() StageRepo {
	return &stageRepo{}
}

func (r *stageRepo) Insert(ctx context.Context, e *model.Stage) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *stageRepo) FindByID(ctx context.Context, id string) (*model.Stage, error) {
	e := &model.Stage{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("cockroach.FindOne: %w", err)
	}

	return e, nil
}

func (r *stageRepo) Update(ctx context.Context, e *model.Stage) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *stageRepo) SoftDelete(ctx context.Context, id string) error {
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
	IDs        []string
	Title      string
	SyllabusID string
	CourseID   string
	Limit      int64
	Offset     int64
}

func (s *SearchStagesOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.StageFieldID)
	}

	if s.Title != "" {
		args = append(args, "%"+s.Title+"%")
		conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.StageFieldTitle, len(args))
	}
	if s.SyllabusID != "" {
		args = append(args, s.SyllabusID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.StageFieldSyllabusID, len(args))
	}
	if s.CourseID != "" {
		args = append(args, s.CourseID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.StageFieldCourseID, len(args))
	}

	b := &model.Stage{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf(`SELECT count(*) as cnt
		FROM %s AS b %s
		WHERE TRUE %s AND b.deleted_at IS NULL`, b.TableName(), joins, conds), args
	}
	return fmt.Sprintf(`SELECT b.%s
		FROM %s AS b %s
		WHERE TRUE %s AND b.deleted_at IS NULL
		ORDER BY b.id DESC
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", b."), b.TableName(), joins, conds, s.Limit, s.Offset), args
}

func (r *stageRepo) Search(ctx context.Context, s *SearchStagesOpts) ([]*model.Stage, error) {
	stages := make([]*model.Stage, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&stages)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return stages, nil
}
func (r *stageRepo) Count(ctx context.Context, s *SearchStagesOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("stage.Count: %w", err)
	}

	return countResult, nil
}
