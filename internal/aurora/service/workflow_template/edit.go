package workflow_template

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

func (c *workflowTemplateService) EditWorkflowTemplate(ctx context.Context, opt *EditWorkflowTemplateOpts) error {
	workflowTemplate := model.WorkflowTemplate{
		ID:        opt.ID,
		Name:      opt.Name,
		Config:    opt.ConfigData,
		Status:    opt.Status,
		UpdatedBy: opt.UpdatedBy,
		UpdatedAt: time.Now(),
	}

	errTx := cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		err := c.workflowTemplateRepo.Update(ctx2, &workflowTemplate)
		if err != nil {
			return fmt.Errorf("c.workflowTemplateRepo.Update: %w", err)
		}

		return nil
	})
	if errTx != nil {
		return errTx
	}

	return nil
}

type EditWorkflowTemplateOpts struct {
	ID         string
	Name       string
	ConfigData any
	Status     enum.CommonStatus
	UpdatedBy  string
}
