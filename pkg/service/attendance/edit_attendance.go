package attendance

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

func (b *attendanceService) EditAttendance(ctx context.Context, opt *EditAttendanceOpts) error {
	_, err := b.FindAttendanceByID(ctx, opt.ID)
	if err != nil {
		return fmt.Errorf("find attendance not found %s: %w", opt.ID, err)
	}

	errTx := cockroach.ExecInTx(ctx, func(ctx context.Context) error {
		var err error
		table := model.Attendance{}
		updater := cockroach.NewUpdater(table.TableName(), model.AttendanceFieldID, opt.ID)
		updater.Set(model.AttendanceFieldCourseID, opt.CourseID)
		updater.Set(model.AttendanceFieldStageID, opt.StageID)
		updater.Set(model.AttendanceFieldLessonID, opt.LessonID)
		updater.Set(model.AttendanceFieldNote, opt.Note)
		updater.Set(model.AttendanceFieldScoreFactor, opt.ScoreFactor)
		updater.Set(model.AttendanceFieldUpdatedAt, time.Now())
		updater.Set(model.AttendanceFieldUpdatedBy, opt.UserID)

		err = cockroach.UpdateFields(ctx, updater)
		if err != nil {
			return fmt.Errorf("update attendance failed %w", err)
		}
		for _, detail := range opt.Detail {
			if detail.ID == "" {
				err = b.attendanceDetailRepo.Insert(ctx, &model.AttendanceDetail{
					ID:           idutil.ULIDNow(),
					StudentID:    detail.StudentID,
					AttendanceID: opt.ID,
					Status:       detail.Status,
					Point:        detail.Point,
				})
				if err != nil {
					return err
				}
			} else {
				detailTable := model.AttendanceDetail{}
				updater := cockroach.NewUpdater(detailTable.TableName(), model.AttendanceDetailFieldID, detail.ID)
				updater.Set(model.AttendanceDetailFieldStatus, detail.Status)
				updater.Set(model.AttendanceDetailFieldStudentID, detail.StudentID)
				updater.Set(model.AttendanceDetailFieldPoint, detail.Point)
				err = cockroach.UpdateFields(ctx, updater)
				if err != nil {
					return fmt.Errorf("update attendance failed %w", err)
				}
			}
		}
		return nil
	})
	if errTx != nil {
		return errTx
	}
	return nil
}

type EditAttendanceOpts struct {
	ID          string
	CourseID    string
	StageID     string
	LessonID    string
	Note        string
	ScoreFactor enum.ScoreFactor
	Detail      []*Detail
	UserID      string
}
