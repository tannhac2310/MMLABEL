package department

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

func (c *departmentService) EditDepartment(ctx context.Context, opt *EditDepartmentOpts) error {
	var err error
	table := model.Department{}
	updater := cockroach.NewUpdater(table.TableName(), model.DepartmentFieldID, opt.ID)

	updater.Set(model.DepartmentFieldName, opt.Name)
	updater.Set(model.DepartmentFieldParentID, opt.ParentID)
	updater.Set(model.DepartmentFieldName, opt.Name)
	updater.Set(model.DepartmentFieldShortName, opt.ShortName)
	updater.Set(model.DepartmentFieldCode, opt.Code)
	updater.Set(model.DepartmentFieldPriority, opt.Priority)

	updater.Set(model.DepartmentFieldUpdatedAt, time.Now())

	err = cockroach.UpdateFields(ctx, updater)
	if err != nil {
		return fmt.Errorf("update department failed %w", err)
	}
	return nil
}

type EditDepartmentOpts struct {
	ID        string
	ParentID  string
	Name      string
	ShortName string
	Code      string
	Priority  int64
}
