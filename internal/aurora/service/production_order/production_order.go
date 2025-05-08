package production_order

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/role"

	"github.com/go-redis/redis"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/configs"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	repository2 "mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

type Service interface {
	CreateProductionOrderStage(ctx context.Context, poId string, opt *ProductionOrderStage) (string, error)
	EditProductionOrderStage(ctx context.Context, opt *ProductionOrderStage) error
	DeleteProductionOrderStage(ctx context.Context, id string) error
	AcceptAndChangeNextStage(ctx context.Context, id string) error
	CreateProductionOrder(ctx context.Context, opt *CreateProductionOrderOpts) (string, error)
	EditProductionOrder(ctx context.Context, opt *EditProductionOrderOpts) error
	FindProductionOrders(ctx context.Context, opts *FindProductionOrdersOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error)
	FindAnalysis(ctx context.Context, opts *FindProductionOrdersOpts) ([]*Analysis, error)
	FindProductionOrdersWithNoPermission(ctx context.Context, opts *FindProductionOrdersOpts, sort *repository.Sort, limit, offset int64) ([]*DataWithNoPermission, *repository.CountResult, error)
	DeleteProductionOrder(ctx context.Context, id string) error
	GetCustomField() []string
}

type productionOrderService struct {
	productionOrderRepo            repository.ProductionOrderRepo
	productionOrderStageRepo       repository.ProductionOrderStageRepo
	productionOrderStageDeviceRepo repository.ProductionOrderStageDeviceRepo
	productionPlanRepo             repository.ProductionPlanRepo
	customFieldRepo                repository.CustomFieldRepo
	customerRepo                   repository.CustomerRepo
	userRepo                       repository2.UserRepo
	orderRepo                      repository.OrderRepo
	stageRepo                      repository.StageRepo
	roleService                    role.Service
	cfg                            *configs.Config
	redisDB                        redis.Cmdable
}

func (c *productionOrderService) GetCustomField() []string {
	return []string{
		// "ma_xuat_film",
		// "ma_dao_be_keo",
		// "ma_dao_thanh_pham",
		// "ma_dao_khuon_dap",
		// "ma_dao_khuon_khiem_thi",
		// "khoi_luong_thanh_pham",
		// "so_mau_in",
		// "so_lan_in",
		// "vat_lieu_chinh",
		// "keo",
		// "mang",
		// "khac",
		// "kho_in",
		// "so_sp_in",
		// "hinh_thuc_in",
		// "so_luong_su_dung",
		// "so_luong_san_xuat",
		// "kich_thuoc_thanh_pham",
		"ma_sp_kh",
		"ma_sp_mm",
		"ma_xuat_phim",
		"ma_dao_be_keo",
		"ma_dao_thanh_pham",
		"ma_khuon_dap",
		"ma_khuon_khiem_thi",
		"kich_thuoc_thanh_pham",
		"so_mau_in",
		"so_lan_in",
		"vat_lieu_chinh",
		"keo",
		"mang",
		"khac",
		"kho_in",
		"so_sp_in_ban_in",
		"hinh_thuc_in",
		"so_luong_su_dung",
		"su_co",
	}
}
func (c *productionOrderService) deleteProductionOrderStage(ctx context.Context, ids []string, productionId string) interface{} {
	// find production order stage by production order id
	productionOrderStages, err := c.productionOrderStageRepo.Search(ctx, &repository.SearchProductionOrderStagesOpts{
		ProductionOrderID: productionId,
		Limit:             10000,
	})
	if err != nil {
		return err
	}
	// get production order stage id to delete
	// find ids not in productionOrderStages.id
	var idsToDelete []string
	for _, id := range ids {
		var found bool
		for _, productionOrderStage := range productionOrderStages {
			if id == productionOrderStage.ID {
				found = true
				break
			}
		}
		if !found {
			idsToDelete = append(idsToDelete, id)
		}
	}
	if len(idsToDelete) <= 0 {
		return nil
	}
	return c.productionOrderStageRepo.SoftDeletes(ctx, idsToDelete)
}

func NewService(
	productionOrderRepo repository.ProductionOrderRepo,
	productionOrderStageRepo repository.ProductionOrderStageRepo,
	productOrderStageDeviceRepo repository.ProductionOrderStageDeviceRepo,
	customFieldRepo repository.CustomFieldRepo,
	customerRepo repository.CustomerRepo,
	userRepo repository2.UserRepo,
	stageRepo repository.StageRepo,
	productionPlanRepo repository.ProductionPlanRepo,
	orderRepo repository.OrderRepo,
	cfg *configs.Config,
	redisDB redis.Cmdable,
	roleService role.Service,
) Service {
	return &productionOrderService{
		productionOrderRepo:            productionOrderRepo,
		productionOrderStageRepo:       productionOrderStageRepo,
		productionOrderStageDeviceRepo: productOrderStageDeviceRepo,
		customFieldRepo:                customFieldRepo,
		productionPlanRepo:             productionPlanRepo,
		customerRepo:                   customerRepo,
		userRepo:                       userRepo,
		orderRepo:                      orderRepo,
		stageRepo:                      stageRepo,
		cfg:                            cfg,
		redisDB:                        redisDB,
		roleService:                    roleService,
	}
}

type Data struct {
	*repository.ProductionOrderData
	ProductionOrderStage []*ProductionOrderStageData
	ProductionPlanData   *repository.ProductionPlanData
	OrderData            *repository.OrderData
	CustomData           map[string]string
	CustomerData         *repository.CustomerData
	CreatedByName        string
}
type DataWithNoPermission struct {
	*repository.ProductionOrderData
}

type ProductionOrderStageData struct {
	*model.ProductionOrderStage
	StageName                  string
	ProductionOrderStageDevice []*repository.ProductionOrderStageDeviceData
}
