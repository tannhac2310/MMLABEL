package stage

import (
	"context"
	"fmt"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

func (c *stageService) EditStage(ctx context.Context, opt *EditStageOpts) error {
	var err error
	table := model.Stage{}
	updater := cockroach.NewUpdater(table.TableName(), model.StageFieldID, opt.ID)

	updater.Set(model.StageFieldName, opt.Name)
	updater.Set(model.StageFieldParentID, opt.ParentID)
	updater.Set(model.StageFieldDepartmentCode, opt.DepartmentCode)
	updater.Set(model.StageFieldName, opt.Name)
	updater.Set(model.StageFieldShortName, opt.ShortName)
	updater.Set(model.StageFieldCode, opt.Code)
	updater.Set(model.StageFieldSorting, opt.Sorting)
	updater.Set(model.StageFieldErrorCode, opt.ErrorCode)
	updater.Set(model.StageFieldData, opt.Data)
	updater.Set(model.StageFieldNote, opt.Note)
	updater.Set(model.StageFieldStatus, opt.Status)

	updater.Set(model.StageFieldUpdatedAt, time.Now())

	err = cockroach.UpdateFields(ctx, updater)
	if err != nil {
		return fmt.Errorf("update stage failed %w", err)
	}
	return nil
}

type EditStageOpts struct {
	ID             string
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
}
