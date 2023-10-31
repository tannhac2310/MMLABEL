package stagestudent

import (
	"context"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

type Service interface {
	CreateStageStudents(ctx context.Context, opt *CreateStageStudentOpts) error
	EditStageStudent(ctx context.Context, opt *EditStageStudentOpts) error
	FindStageStudents(ctx context.Context, opts *FindStageStudentsOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error)
	SoftDelete(ctx context.Context, id string) error
	FindStageStudentByID(ctx context.Context, id string) (*Data, error)
}

type stageStudentService struct {
	stageStudentRepo     repository.StageStudentRepo
	attendanceDetailRepo repository.AttendanceDetailRepo
}

func NewService(
	stageStudentRepo repository.StageStudentRepo,
	attendanceDetailRepo repository.AttendanceDetailRepo,
) Service {
	return &stageStudentService{
		stageStudentRepo:     stageStudentRepo,
		attendanceDetailRepo: attendanceDetailRepo,
	}
}

type Attendance struct {
	ID         string            `db:"id"`
	StudentID  string            `db:"student_id"`
	CourseID   string            `db:"course_id"`
	StageID    string            `db:"stage_id"`
	LessonID   string            `db:"lesson_id"`
	Status     enum.CommonStatus `db:"status"`
	Point      float64           `db:"point"`
	RecordedAt time.Time         `db:"recorded_at"`
}
type Data struct {
	*repository.StageStudentData
	Detail []*repository.AttendanceDetailInfo
}
