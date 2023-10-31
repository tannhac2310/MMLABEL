package report

import (
	"context"

	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

type StudentAttendance struct {
	StageID          string
	StageName        string
	AttendanceStatus enum.AttendanceDetailStatus
	Date             time.Time
	Count            int64
}
type OrderByDate struct {
	Date  time.Time
	Total float64
}
type Data struct {
	StageCount     int64
	StudentInStage []*repository.StudentInStageData
	OrderByDate    []*repository.OrderByDate
}
type Service interface {
	DashboardReport(ctx context.Context) (*Data, error)
}
type reportService struct {
	courseRepo       repository.CourseRepo
	stageRepo        repository.StageRepo
	stageStudentRepo repository.StageStudentRepo
	orderRepo        repository.OrderRepo
}

func (r *reportService) DashboardReport(ctx context.Context) (*Data, error) {
	total, err := r.stageRepo.Count(ctx, &repository.SearchStagesOpts{})
	if err != nil {
		return nil, err
	}
	studentInStage, errOrder := r.stageStudentRepo.ReportStudentInStage(ctx)
	if errOrder != nil {
		return nil, errOrder
	}
	orderByDate, errOrderByDate := r.orderRepo.ReportOrderByDate(ctx)
	if errOrderByDate != nil {
		return nil, errOrderByDate
	}
	return &Data{
		StageCount:     total.Count,
		StudentInStage: studentInStage,
		OrderByDate:    orderByDate,
	}, nil
}

func NewService(
	courseRepo repository.CourseRepo,
	stageRepo repository.StageRepo,
	stageStudentReport repository.StageStudentRepo,
	orderRepo repository.OrderRepo,
) Service {
	return &reportService{
		courseRepo:       courseRepo,
		stageRepo:        stageRepo,
		stageStudentRepo: stageStudentReport,
		orderRepo:        orderRepo,
	}
}
