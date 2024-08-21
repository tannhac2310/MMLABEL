package workflow_template

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

func (c *workflowTemplateService) FindWorkflowTemplates(ctx context.Context, opts *FindWorkflowTemplatesOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error) {
	filter := &repository.SearchWorkflowTemplateOpts{
		IDs:    opts.IDs,
		Name:   opts.Name,
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	}
	workflowTemplates, err := c.workflowTemplateRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	total, err := c.workflowTemplateRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*Data, 0, len(workflowTemplates))
	for _, workflowTemplate := range workflowTemplates {
		results = append(results, &Data{
			WorkflowTemplateData: workflowTemplate,
		})
	}
	return results, total, nil
}

type FindWorkflowTemplatesOpts struct {
	IDs   []string
	Name  string
	Phone string
	Code  string
}
