package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

type AttendanceRepo interface {
	Insert(ctx context.Context, e *model.Attendance) error
	Update(ctx context.Context, e *model.Attendance) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchAttendancesOpts) ([]*AttendanceData, error)
	Count(ctx context.Context, s *SearchAttendancesOpts) (*CountResult, error)
}

type attendancesRepo struct {
}

func NewAttendanceRepo() AttendanceRepo {
	return &attendancesRepo{}
}

func (r *attendancesRepo) Insert(ctx context.Context, e *model.Attendance) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *attendancesRepo) Update(ctx context.Context, e *model.Attendance) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *attendancesRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE attendance
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

// SearchAttendancesOpts all params is options
type SearchAttendancesOpts struct {
	IDs      []string
	CourseID string
	StageID  string
	LessonID string
	Limit    int64
	Offset   int64
	Sort     *Sort
}

func (s *SearchAttendancesOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.AttendanceFieldID)
	}

	if s.CourseID != "" {
		args = append(args, s.CourseID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.AttendanceFieldCourseID, len(args))
	}

	if s.StageID != "" {
		args = append(args, s.StageID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.AttendanceFieldStageID, len(args))
	}
	if s.LessonID != "" {
		args = append(args, s.LessonID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.AttendanceFieldLessonID, len(args))
	}

	b := &model.Attendance{}
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
	return fmt.Sprintf(`SELECT b.%s, c.title as course_name, s.title as stage_name, l.title as lesson_name, cb.name as created_by_name, ub.name as updated_by_name 
		FROM %s AS b %s
		JOIN courses c on c.id = b.course_id
		JOIN stages s on s.id = b.stage_id
		JOIN lessons l on l.id = b.lesson_id
		JOIN users cb on cb.id = b.created_by
		JOIN users ub on ub.id = b.updated_by
		WHERE TRUE %s AND b.deleted_at IS NULL
		%s
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type AttendanceData struct {
	*model.Attendance
	CourseName    string `db:"course_name"`
	StageName     string `db:"stage_name"`
	LessonName    string `db:"lesson_name"`
	CreatedByName string `db:"created_by_name"`
	UpdatedByName string `db:"updated_by_name"`
}

func (r *attendancesRepo) Search(ctx context.Context, s *SearchAttendancesOpts) ([]*AttendanceData, error) {
	attendance := make([]*AttendanceData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&attendance)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return attendance, nil
}

func (r *attendancesRepo) Count(ctx context.Context, s *SearchAttendancesOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("attendances.Count: %w", err)
	}

	return countResult, nil
}
