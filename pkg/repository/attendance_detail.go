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

type AttendanceDetailRepo interface {
	Insert(ctx context.Context, e *model.AttendanceDetail) error
	Update(ctx context.Context, e *model.AttendanceDetail) error
	FindByID(ctx context.Context, id string) (*model.AttendanceDetail, error)
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchAttendanceDetailsOpts) ([]*AttendanceDetailInfo, error)
	Count(ctx context.Context, s *SearchAttendanceDetailsOpts) (*CountResult, error)
}

type attendanceDetailsRepo struct {
}

func NewAttendanceDetailRepo() AttendanceDetailRepo {
	return &attendanceDetailsRepo{}
}

func (r *attendanceDetailsRepo) Insert(ctx context.Context, e *model.AttendanceDetail) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *attendanceDetailsRepo) FindByID(ctx context.Context, id string) (*model.AttendanceDetail, error) {
	e := &model.AttendanceDetail{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("cockroach.FindOne: %w", err)
	}

	return e, nil
}

func (r *attendanceDetailsRepo) Update(ctx context.Context, e *model.AttendanceDetail) error {
	return cockroach.Update(ctx, e)
}

func (r *attendanceDetailsRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE attendanceDetail
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

// SearchAttendanceDetailsOpts all params is options
type SearchAttendanceDetailsOpts struct {
	IDs          []string
	AttendanceID string
	StudentID    string
	Status       enum.AttendanceDetailStatus
	LessonID     string
	Limit        int64
	Offset       int64
	Sort         *Sort
}

func (s *SearchAttendanceDetailsOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.AttendanceDetailFieldID)
	}

	if s.AttendanceID != "" {
		args = append(args, s.AttendanceID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.AttendanceDetailFieldAttendanceID, len(args))
	}

	if s.StudentID != "" {
		args = append(args, s.StudentID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.AttendanceDetailFieldStudentID, len(args))
	}
	if s.Status > 0 {
		args = append(args, s.Status)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.AttendanceDetailFieldStatus, len(args))
	}

	b := &model.AttendanceDetail{}
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

	return fmt.Sprintf(`SELECT b.%s, u.name as student_name, u.email as student_email, u.avatar as student_avatar,
				a.lesson_id as lesson_id, l.title as lesson_name, a.recorded_at as recorded_at, a.score_factor as score_factor
		FROM %s AS b %s
		RIGHT JOIN attendances a on a.id = b.attendance_id
		RIGHT JOIN lessons l on l.id = a.lesson_id
		JOIN users u on u.id = b.student_id
		WHERE TRUE %s AND b.deleted_at IS NULL
		%s
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type AttendanceDetailInfo struct {
	*model.AttendanceDetail
	RecordedAt    time.Time        `db:"recorded_at"`
	ScoreFactor   enum.ScoreFactor `db:"score_factor"`
	LessonID      string           `db:"lesson_id"`
	LessonName    string           `db:"lesson_name"`
	StudentName   string           `db:"student_name"`
	StudentEmail  string           `db:"student_email"`
	StudentAvatar string           `db:"student_avatar"`
}

func (r *attendanceDetailsRepo) Search(ctx context.Context, s *SearchAttendanceDetailsOpts) ([]*AttendanceDetailInfo, error) {
	attendanceDetail := make([]*AttendanceDetailInfo, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&attendanceDetail)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return attendanceDetail, nil
}

func (r *attendanceDetailsRepo) Count(ctx context.Context, s *SearchAttendanceDetailsOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	fmt.Println(sql)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("attendanceDetails.Count: %w", err)
	}

	return countResult, nil
}
