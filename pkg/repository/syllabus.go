package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

type SyllabusRepo interface {
	Insert(ctx context.Context, e *model.Syllabus) error
	Update(ctx context.Context, e *model.Syllabus) error
	FindByID(ctx context.Context, id string) (*model.Syllabus, error)
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchSyllabusesOpts) ([]*model.Syllabus, error)
	Count(ctx context.Context, s *SearchSyllabusesOpts) (*CountResult, error)
}

type syllabusRepo struct {
}

func NewSyllabusRepo() SyllabusRepo {
	return &syllabusRepo{}
}

func (r *syllabusRepo) Insert(ctx context.Context, e *model.Syllabus) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *syllabusRepo) FindByID(ctx context.Context, id string) (*model.Syllabus, error) {
	e := &model.Syllabus{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("cockroach.FindOne: %w", err)
	}

	return e, nil
}

func (r *syllabusRepo) Update(ctx context.Context, e *model.Syllabus) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *syllabusRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE syllabuses
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
type SearchSyllabusesOpts struct {
	IDs       []string
	Search    string
	TeacherID string
	CourseID  string
	Limit     int64
	Offset    int64
	Sort      *Sort
}

func (s *SearchSyllabusesOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.SyllabusFieldID)
	}

	if s.Search != "" {
		args = append(args, "%"+s.Search+"%")
		// Title
		conds += fmt.Sprintf(" AND (b.%s ILIKE $%d", model.SyllabusFieldTitle, len(args))
		// Code
		conds += fmt.Sprintf(" OR b.%s ILIKE $%d )", model.SyllabusFieldCode, len(args))
	}
	if s.TeacherID != "" {
		args = append(args, s.TeacherID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.StageFieldTeacherID, len(args))
	}
	if s.CourseID != "" {
		args = append(args, s.CourseID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.StageFieldCourseID, len(args))
	}

	b := &model.Syllabus{}
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

func (r *syllabusRepo) Search(ctx context.Context, s *SearchSyllabusesOpts) ([]*model.Syllabus, error) {
	syllabuses := make([]*model.Syllabus, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&syllabuses)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return syllabuses, nil
}

func (r *syllabusRepo) Count(ctx context.Context, s *SearchSyllabusesOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	fmt.Println(sql)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("syllabus.Count: %w", err)
	}

	return countResult, nil
}
