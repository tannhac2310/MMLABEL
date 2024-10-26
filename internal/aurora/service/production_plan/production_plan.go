package production_plan

import (
	"context"

	"github.com/go-redis/redis"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/configs"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	repository2 "mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/role"
)

type Service interface {
	FindProductionPlans(ctx context.Context, opts *FindProductionPlansOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error)
	FindProductionPlansWithNoPermission(ctx context.Context, opts *FindProductionPlansOpts, sort *repository.Sort, limit, offset int64) ([]*DataWithNoPermission, *repository.CountResult, error)
	CreateProductionPlan(ctx context.Context, opt *CreateProductionPlanOpts) (string, error)
	EditProductionPlan(ctx context.Context, opt *EditProductionPlanOpts) error
	DeleteProductionPlan(ctx context.Context, id string) error
	ProcessProductionOrder(ctx context.Context, opt *ProcessProductionOrderOpts) (string, error)
	UpdateCustomFields(ctx context.Context, productionPlanID string, values []*CustomField) error
	UpdateCurrentStage(ctx context.Context, productionPlanID string, stageID enum.ProductionPlanStage) error
	SummaryProductionPlans(ctx context.Context, opts *SummaryProductionPlanOpts) ([]*SummaryData, error)
}

type productionPlanService struct {
	productionPlanRepo       repository.ProductionPlanRepo
	productionOrderRepo      repository.ProductionOrderRepo
	productionOrderStageRepo repository.ProductionOrderStageRepo
	customFieldRepo          repository.CustomFieldRepo
	customerRepo             repository.CustomerRepo
	userRepo                 repository2.UserRepo
	deviceConfigRepo         repository.ProductionOrderDeviceConfigRepo
	roleService              role.Service
	cfg                      *configs.Config
	redisDB                  redis.Cmdable
}

func NewService(
	productionPlanRepo repository.ProductionPlanRepo,
	productionOrderRepo repository.ProductionOrderRepo,
	productionOrderStageRepo repository.ProductionOrderStageRepo,
	customFieldRepo repository.CustomFieldRepo,
	customerRepo repository.CustomerRepo,
	userRepo repository2.UserRepo,
	roleService role.Service,
	deviceConfigRepo repository.ProductionOrderDeviceConfigRepo,
	cfg *configs.Config,
	redisDB redis.Cmdable,
) Service {
	return &productionPlanService{
		productionPlanRepo:       productionPlanRepo,
		productionOrderRepo:      productionOrderRepo,
		productionOrderStageRepo: productionOrderStageRepo,
		customFieldRepo:          customFieldRepo,
		customerRepo:             customerRepo,
		userRepo:                 userRepo,
		deviceConfigRepo:         deviceConfigRepo,
		roleService:              roleService,
		cfg:                      cfg,
		redisDB:                  redisDB,
	}
}

type Data struct {
	*repository.ProductionPlanData
	CustomData   map[string]string
	UserFields   map[string][]*repository.CustomFieldData
	CustomerData *repository.CustomerData
}

type DataWithNoPermission struct {
	*repository.ProductionPlanData
	CustomerData *repository.CustomerData
}

type SummaryData struct {
	*repository.SummaryProductionPlanData
}
