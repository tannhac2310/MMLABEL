package stage

import (
	"context"
	"database/sql"
	"fmt"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

func (c *stageService) CreateStage(ctx context.Context, opt *CreateStageOpts) (string, error) {
	now := time.Now()

	stage := &model.Stage{
		ID:             idutil.ULIDNow(),
		ParentID:       cockroach.String(opt.ParentID),
		DepartmentCode: cockroach.String(opt.DepartmentCode),
		Name:           opt.Name,
		ShortName:      opt.ShortName,
		Code:           opt.Code,
		Sorting:        opt.Sorting,
		ErrorCode:      cockroach.String(opt.ErrorCode),
		Data:           opt.Data,
		Note:           cockroach.String(opt.Note),
		Status:         opt.Status,
		CreatedBy:      opt.CreatedBy,
		CreatedAt:      now,
		UpdatedAt:      now,
		DeletedAt:      sql.NullTime{},
	}

	errTx := cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		err := c.stageRepo.Insert(ctx2, stage)
		if err != nil {
			return fmt.Errorf("c.stageRepo.Insert: %w", err)
		}

		return nil
	})
	if errTx != nil {
		return "", errTx
	}
	return stage.ID, nil
}

type CreateStageOpts struct {
	ParentID       string
	DepartmentCode string
	Name           string
	ShortName      string
	Code           string
	Sorting        int16
	ErrorCode      string
	Data           map[string]interface{}
	Note           string
	Status         enum.StageStatus
	CreatedBy      string
}
