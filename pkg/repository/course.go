package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

type CourseRepo interface {
	Insert(ctx context.Context, e *model.Course) error
	Update(ctx context.Context, e *model.Course) error
	FindByID(ctx context.Context, id string) (*model.Course, error)
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchCoursesOpts) ([]*model.Course, error)
	Count(ctx context.Context, s *SearchCoursesOpts) (*CountResult, error)
}

type courseRepo struct {
}

func NewCourseRepo() CourseRepo {
	return &courseRepo{}
}

func (r *courseRepo) Insert(ctx context.Context, e *model.Course) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *courseRepo) FindByID(ctx context.Context, id string) (*model.Course, error) {
	e := &model.Course{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("CourseRepo:cockroach.FindOne: %w, id=%s", err, id)
	}

	return e, nil
}

func (r *courseRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE courses
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

func (r *courseRepo) Update(ctx context.Context, e *model.Course) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

// SearchCoursesOpts all params is options
type SearchCoursesOpts struct {
	IDs    []string
	Code   string
	Search string
	Title  string
	Type   enum.CourseType
	Status enum.CommonStatus
	Limit  int64
	Offset int64
}

func (s *SearchCoursesOpts) buildQuery(isCount bool) (string, []interface{}) {
	args := []interface{}{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND c.%s = ANY($1)", model.CourseFieldID)
	}

	if s.Title != "" {
		args = append(args, "%"+s.Title+"%")
		conds += fmt.Sprintf(" AND c.%s ILIKE $%d", model.CourseFieldTitle, len(args))
	}
	if s.Search != "" {
		// Title
		args = append(args, "%"+s.Search+"%")
		conds += fmt.Sprintf(" AND (c.%s ILIKE $%d", model.CourseFieldTitle, len(args))
		// Code
		conds += fmt.Sprintf(" OR c.%s ILIKE $%d", model.CourseFieldCode, len(args))
		// Code
		conds += fmt.Sprintf(" OR c.%s ILIKE $%d )", model.CourseFieldCode, len(args))
	}

	if s.Code != "" {
		args = append(args, s.Code)
		conds += fmt.Sprintf(" AND c.%s = $%d", model.CourseFieldCode, len(args))
	}

	if s.Title != "" {
		args = append(args, s.Title)
		conds += fmt.Sprintf(" AND c.%s = $%d", model.CourseFieldTitle, len(args))
	}

	if s.Type > 0 {
		args = append(args, s.Type)
		conds += fmt.Sprintf(" AND c.%s = $%d", model.CourseFieldType, len(args))
	}

	if s.Status > 0 {
		args = append(args, s.Status)
		conds += fmt.Sprintf(" AND c.%s = $%d", model.CourseFieldStatus, len(args))
	}

	c := &model.Course{}
	fields, _ := c.FieldMap()
	if isCount == true {
		return fmt.Sprintf(`SELECT count(*) as cnt
		FROM %s AS c %s
		WHERE TRUE %s AND c.deleted_at IS NULL`, c.TableName(), joins, conds), args
	}
	return fmt.Sprintf(`SELECT c.%s
		FROM %s AS c %s
		WHERE TRUE %s AND c.deleted_at IS NULL
		ORDER BY c.id DESC
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", c."), c.TableName(), joins, conds, s.Limit, s.Offset), args
}

func (r *courseRepo) Search(ctx context.Context, s *SearchCoursesOpts) ([]*model.Course, error) {
	courses := make([]*model.Course, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&courses)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return courses, nil
}
func (r *courseRepo) Count(ctx context.Context, s *SearchCoursesOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	fmt.Println(sql)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("course.Count: %w", err)
	}

	return countResult, nil
}
