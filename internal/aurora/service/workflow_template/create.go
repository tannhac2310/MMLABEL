package workflow_template

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

func (c *workflowTemplateService) CreateWorkflowTemplate(ctx context.Context, opt *CreateWorkflowTemplateOpts) (string, error) {
	now := time.Now()

	workflowTemplate := &model.WorkflowTemplate{
		ID:        idutil.ULIDNow(),
		Name:      opt.Name,
		Config:    opt.ConfigData,
		Status:    opt.Status,
		CreatedBy: opt.CreatedBy,
		UpdatedBy: opt.CreatedBy,
		CreatedAt: now,
		UpdatedAt: now,
	}
	errTx := cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		err := c.workflowTemplateRepo.Insert(ctx2, workflowTemplate)
		if err != nil {
			return fmt.Errorf("c.workflowTemplateRepo.Insert: %w", err)
		}

		return nil
	})
	if errTx != nil {
		return "", errTx
	}

	return workflowTemplate.ID, nil
}

type CreateWorkflowTemplateOpts struct {
	Name       string
	ConfigData any
	Status     enum.CommonStatus
	CreatedBy  string
}
