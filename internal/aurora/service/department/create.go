package department

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

func (c *departmentService) CreateDepartment(ctx context.Context, opt *CreateDepartmentOpts) (string, error) {
	now := time.Now()

	department := &model.Department{
		ID:        idutil.ULIDNow(),
		ParentID:  cockroach.String(opt.ParentID),
		Name:      opt.Name,
		ShortName: opt.ShortName,
		Code:      opt.Code,
		Priority:  opt.Priority,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: sql.NullTime{},
	}

	errTx := cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		err := c.departmentRepo.Insert(ctx2, department)
		if err != nil {
			return fmt.Errorf("c.departmentRepo.Insert: %w", err)
		}

		return nil
	})
	if errTx != nil {
		return "", errTx
	}
	return department.ID, nil
}

type CreateDepartmentOpts struct {
	ParentID  string
	Name      string
	ShortName string
	Code      string
	Priority  int64
	CreatedBy string
}
