package stagestudent

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (b stageStudentService) CreateStageStudents(ctx context.Context, opt *CreateStageStudentOpts) error {
	now := time.Now()
	errTx := cockroach.ExecInTx(ctx, func(ctx context.Context) error {
		for _, id := range opt.StudentIDs {
			stageStudent := &model.StageStudent{
				ID:        idutil.ULIDNow(),
				StudentID: id,
				CourseID:  opt.CourseID,
				StageID:   opt.StageID,
				Status:    opt.Status,
				CreatedBy: opt.CreatedBy,
				UpdatedBy: opt.CreatedBy,
				CreatedAt: now,
				UpdatedAt: now,
			}
			err := b.stageStudentRepo.Insert(ctx, stageStudent)
			if err != nil {
				return fmt.Errorf("%s,p.stageStudentRepo.Insert: %w", "100", err)
			}
		}
		return nil
	})
	if errTx != nil {
		return fmt.Errorf("p.stageStudentRepo.Insert: %w", errTx)
	}
	return nil
}

type CreateStageStudentOpts struct {
	StudentIDs []string
	CourseID   string
	StageID    string
	Status     enum.StageStudentStatus
	CreatedBy  string
}
