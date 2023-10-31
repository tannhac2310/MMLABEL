package attendance

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

type Service interface {
	CreateAttendance(ctx context.Context, opt *CreateAttendanceOpts) (string, error)
	EditAttendance(ctx context.Context, opt *EditAttendanceOpts) error
	FindAttendances(ctx context.Context, opts *FindAttendancesOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error)
	SoftDelete(ctx context.Context, id string) error
	FindAttendanceByID(ctx context.Context, id string) (*Data, error)
}

type attendanceService struct {
	attendanceRepo       repository.AttendanceRepo
	attendanceDetailRepo repository.AttendanceDetailRepo
	userRepo             repository.UserRepo
	courseRepo           repository.CourseRepo
}

func NewService(
	attendanceRepo repository.AttendanceRepo,
	userRepo repository.UserRepo,
	courseRepo repository.CourseRepo,
	attendanceDetailRepo repository.AttendanceDetailRepo,
) Service {
	return &attendanceService{
		attendanceRepo:       attendanceRepo,
		userRepo:             userRepo,
		courseRepo:           courseRepo,
		attendanceDetailRepo: attendanceDetailRepo,
	}
}

type Data struct {
	*repository.AttendanceData
	Detail        []*repository.AttendanceDetailInfo
	CreatedByInfo *model.UserInfo
	UpdatedByInfo *model.UserInfo
}
