package production_plan

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/role"

	"github.com/go-redis/redis"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/configs"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	repository2 "mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

type Service interface {
	FindProductionPlans(ctx context.Context, opts *FindProductionPlansOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error)
	FindProductionPlansWithNoPermission(ctx context.Context, opts *FindProductionPlansOpts, sort *repository.Sort, limit, offset int64) ([]*DataWithNoPermission, *repository.CountResult, error)
	CreateProductionPlan(ctx context.Context, opt *CreateProductionPlanOpts) (string, error)
	EditProductionPlan(ctx context.Context, opt *EditProductionPlanOpts) error
	DeleteProductionPlan(ctx context.Context, id string) error
	GetCustomField() []string
}

type productionPlanService struct {
	productionPlanRepo repository.ProductionPlanRepo
	customFieldRepo    repository.CustomFieldRepo
	userRepo           repository2.UserRepo
	roleService        role.Service
	cfg                *configs.Config
	redisDB            redis.Cmdable
}

func NewService(
	productionPlanRepo repository.ProductionPlanRepo,
	customFieldRepo repository.CustomFieldRepo,
	userRepo repository2.UserRepo,
	roleService role.Service,
	cfg *configs.Config,
	redisDB redis.Cmdable,
) Service {
	return &productionPlanService{
		productionPlanRepo: productionPlanRepo,
		customFieldRepo:    customFieldRepo,
		userRepo:           userRepo,
		roleService:        roleService,
		cfg:                cfg,
		redisDB:            redisDB,
	}
}

func (c *productionPlanService) GetCustomField() []string {
	// TODO update this
	return []string{
		"a",
		"b",
	}
}

type Data struct {
	*repository.ProductionPlanData
	CustomData map[string]string
}

type DataWithNoPermission struct {
	*repository.ProductionPlanData
}
