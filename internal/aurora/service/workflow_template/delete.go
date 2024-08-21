package workflow_template

import (
	"context"
)

func (c *workflowTemplateService) Delete(ctx context.Context, id string) error {
	return c.workflowTemplateRepo.SoftDelete(ctx, id)
}
