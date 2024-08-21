package workflow_template

import (
	"context"

	"github.com/go-redis/redis"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/configs"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

type Service interface {
	CreateWorkflowTemplate(ctx context.Context, opt *CreateWorkflowTemplateOpts) (string, error)
	EditWorkflowTemplate(ctx context.Context, opt *EditWorkflowTemplateOpts) error
	FindWorkflowTemplates(ctx context.Context, opts *FindWorkflowTemplatesOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error)
	Delete(ctx context.Context, id string) error
}

type workflowTemplateService struct {
	workflowTemplateRepo repository.WorkflowTemplateRepo
	cfg                  *configs.Config
	redisDB              redis.Cmdable
}

func NewService(
	workflowTemplateRepo repository.WorkflowTemplateRepo,
	cfg *configs.Config,
	redisDB redis.Cmdable,
) Service {
	return &workflowTemplateService{
		workflowTemplateRepo: workflowTemplateRepo,
		cfg:                  cfg,
		redisDB:              redisDB,
	}
}

type Data struct {
	*repository.WorkflowTemplateData
}
