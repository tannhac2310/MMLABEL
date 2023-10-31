package attendance

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (b *attendanceService) CreateAttendance(ctx context.Context, opt *CreateAttendanceOpts) (string, error) {
	now := time.Now()
	attendanceID := idutil.ULIDNow()
	errTx := cockroach.ExecInTx(ctx, func(ctx context.Context) error {
		attendance := &model.Attendance{
			ID:          attendanceID,
			CourseID:    opt.CourseID,
			StageID:     opt.StageID,
			LessonID:    opt.LessonID,
			Note:        opt.Note,
			ScoreFactor: opt.ScoreFactor,
			RecordedAt:  opt.RecordedAt,
			CreatedBy:   opt.CreatedBy,
			UpdatedBy:   opt.CreatedBy,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		err := b.attendanceRepo.Insert(ctx, attendance)
		if err != nil {
			return fmt.Errorf("p.attendanceRepo.Insert: %w", err)
		}
		for _, detail := range opt.Detail {
			err = b.attendanceDetailRepo.Insert(ctx, &model.AttendanceDetail{
				ID:           idutil.ULIDNow(),
				StudentID:    detail.StudentID,
				AttendanceID: attendanceID,
				Status:       detail.Status,
				Point:        detail.Point,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	if errTx != nil {
		return "", errTx
	}
	return attendanceID, nil
}

type Detail struct {
	ID           string
	StudentID    string
	AttendanceID string
	Status       enum.AttendanceDetailStatus
	Point        float64
}
type CreateAttendanceOpts struct {
	CourseID    string
	StageID     string
	LessonID    string
	Note        string
	Detail      []*Detail
	ScoreFactor enum.ScoreFactor
	RecordedAt  time.Time
	CreatedBy   string
}
