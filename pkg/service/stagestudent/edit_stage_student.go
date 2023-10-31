package stagestudent

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (b *stageStudentService) EditStageStudent(ctx context.Context, opt *EditStageStudentOpts) error {
	var err error
	table := model.StageStudent{}
	updater := cockroach.NewUpdater(table.TableName(), model.StageStudentFieldID, opt.ID)

	updater.Set(model.StageStudentFieldStudentID, opt.StudentID)
	updater.Set(model.StageStudentFieldCourseID, opt.CourseID)
	updater.Set(model.StageStudentFieldStageID, opt.StageID)
	updater.Set(model.StageStudentFieldStatus, opt.Status)

	err = cockroach.UpdateFields(ctx, updater)
	if err != nil {
		return fmt.Errorf("err update stage_student status: %w", err)
	}

	return nil
}

type EditStageStudentOpts struct {
	ID        string
	StudentID string
	CourseID  string
	StageID   string
	Status    enum.StageStudentStatus
	UserID    string
}
