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

type StageStudentRepo interface {
	Insert(ctx context.Context, e *model.StageStudent) error
	Update(ctx context.Context, e *model.StageStudent) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchStageStudentsOpts) ([]*StageStudentData, error)
	Count(ctx context.Context, s *SearchStageStudentsOpts) (*CountResult, error)
	ReportStudentInStage(ctx context.Context) ([]*StudentInStageData, error)
}
type StudentInStageData struct {
	StageID     string                  `db:"stage_id"`
	StageName   string                  `db:"stage_name"`
	StageStatus enum.StageStudentStatus `db:"stage_status"`
	Count       int64                   `db:"count"`
}

func (r *stageStudentsRepo) ReportStudentInStage(ctx context.Context) ([]*StudentInStageData, error) {
	sql := `select s.id as stage_id, s.title as stage_name, s.status as stage_status, count(*) as count from stage_students ss 
join stages s on s.id = ss.stage_id
group by s.id, s.title, s.status
`
	data := make([]*StudentInStageData, 0)
	err := cockroach.Select(ctx, sql).ScanAll(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type Student struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Avatar      string          `json:"avatar"`
	Address     string          `json:"address"`
	PhoneNumber string          `json:"phone_number"`
	Email       string          `json:"email"`
	Linked      []string        `json:"linked"`
	Status      enum.UserStatus `json:"status"`
	Type        enum.UserType   `json:"type"`
}
type StageStudentData struct {
	*model.StageStudent
	CourseName    string `db:"course_name"`
	StageName     string `db:"stage_name"`
	StudentName   string `db:"student_name"`
	StudentEmail  string `db:"student_email"`
	StudentAvatar string `db:"student_avatar"`
	CreatedByName string `db:"created_by_name"`
	UpdatedByName string `db:"updated_by_name"`
}
type stageStudentsRepo struct {
}

func NewStageStudentRepo() StageStudentRepo {
	return &stageStudentsRepo{}
}

func (r *stageStudentsRepo) Insert(ctx context.Context, e *model.StageStudent) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *stageStudentsRepo) Update(ctx context.Context, e *model.StageStudent) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *stageStudentsRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE state_student
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
type SearchStageStudentsOpts struct {
	IDs       []string
	CourseID  string
	StageID   string
	StudentID string
	Status    enum.StageStudentStatus
	LessonID  string
	Limit     int64
	Offset    int64
	Sort      *Sort
}

func (s *SearchStageStudentsOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.StageStudentFieldID)
	}

	if s.CourseID != "" {
		args = append(args, s.CourseID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.StageStudentFieldCourseID, len(args))
	}

	if s.StageID != "" {
		args = append(args, s.StageID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.StageStudentFieldStageID, len(args))
	}
	if s.StudentID != "" {
		args = append(args, s.StudentID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.StageStudentFieldStudentID, len(args))
	}
	if s.Status > 0 {
		args = append(args, s.Status)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.StageStudentFieldStatus, len(args))
	}

	b := &model.StageStudent{}
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
		c.title as course_name, s.title as stage_name	
		FROM %s AS b %s
		JOIN courses c on c.id = b.course_id
		JOIN stages s on s.id = b.stage_id
		JOIN users u on u.id = b.student_id
		JOIN users cb on cb.id = b.created_by
		JOIN users ub on ub.id = b.updated_by
		WHERE TRUE %s AND b.deleted_at IS NULL
		%s
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

func (r *stageStudentsRepo) Search(ctx context.Context, s *SearchStageStudentsOpts) ([]*StageStudentData, error) {
	stageStudent := make([]*StageStudentData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&stageStudent)
	if err != nil {
		return nil, fmt.Errorf("stageStudent.Search: %w", err)
	}
	return stageStudent, nil
}

func (r *stageStudentsRepo) Count(ctx context.Context, s *SearchStageStudentsOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	fmt.Println(sql)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("stageStudents.Count: %w", err)
	}

	return countResult, nil
}
