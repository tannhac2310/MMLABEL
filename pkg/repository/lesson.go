package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

type LessonRepo interface {
	Insert(ctx context.Context, e *model.Lesson) error
	Update(ctx context.Context, e *model.Lesson) error
	FindByID(ctx context.Context, id string) (*model.Lesson, error)
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchLessonsOpts) ([]*model.Lesson, error)
	Count(ctx context.Context, s *SearchLessonsOpts) (*CountResult, error)
}

type lessonRepo struct {
}

func NewLessonRepo() LessonRepo {
	return &lessonRepo{}
}

func (r *lessonRepo) Insert(ctx context.Context, e *model.Lesson) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *lessonRepo) FindByID(ctx context.Context, id string) (*model.Lesson, error) {
	e := &model.Lesson{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("cockroach.FindOne: %w", err)
	}

	fmt.Println(2)
	return e, nil
}

func (r *lessonRepo) Update(ctx context.Context, e *model.Lesson) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *lessonRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE lessons
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

// all params is options
type SearchLessonsOpts struct {
	IDs        []string
	Search     string
	SyllabusID string
	Limit      int64
	Offset     int64
	Sort       *Sort
}

func (s *SearchLessonsOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.LessonFieldID)
	}

	if s.Search != "" {
		args = append(args, "%"+s.Search+"%")
		// Title
		conds += fmt.Sprintf(" AND (b.%s ILIKE $%d", model.LessonFieldTitle, len(args))
		// Code
		conds += fmt.Sprintf(" OR b.%s ILIKE $%d )", model.LessonFieldLink, len(args))
	}

	if s.SyllabusID != "" {
		args = append(args, s.SyllabusID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.LessonFieldSyllabusID, len(args))
	}

	b := &model.Lesson{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf(`SELECT count(*) as cnt
		FROM %s AS b %s
		WHERE TRUE %s AND b.deleted_at IS NULL`, b.TableName(), joins, conds), args
	}

	order := ""
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

func (r *lessonRepo) Search(ctx context.Context, s *SearchLessonsOpts) ([]*model.Lesson, error) {
	lessons := make([]*model.Lesson, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&lessons)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return lessons, nil
}

func (r *lessonRepo) Count(ctx context.Context, s *SearchLessonsOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	fmt.Println(sql)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("lesson.Count: %w", err)
	}

	return countResult, nil
}
