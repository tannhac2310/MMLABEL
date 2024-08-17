package production_plan

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/role"

	"github.com/go-redis/redis"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/configs"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	repository2 "mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

const (
	ProductionPlanCustomField_ten_sp              = "ten_sp"
	ProductionPlanCustomField_ma_sp               = "ma_sp"
	ProductionPlanCustomField_ma_sp_phu           = "ma_sp_phu"
	ProductionPlanCustomField_dvt                 = "dvt"
	ProductionPlanCustomField_dai                 = "dai"
	ProductionPlanCustomField_rong                = "rong"
	ProductionPlanCustomField_don_vi_dvt          = "don_vi_dvt"
	ProductionPlanCustomField_so_luong_mau        = "so_luong_mau"
	ProductionPlanCustomField_so_lan_in           = "so_lan_in"
	ProductionPlanCustomField_ma_chat_lieu_in     = "ma_chat_lieu_in"
	ProductionPlanCustomField_vat_lieu_thay_the   = "vat_lieu_thay_the"
	ProductionPlanCustomField_keo_2_mat           = "keo_2_mat"
	ProductionPlanCustomField_ma_keo_2_mat        = "ma_keo_2_mat"
	ProductionPlanCustomField_so_luong_keo_2_mat  = "so_luong_keo_2_mat"
	ProductionPlanCustomField_mieu_ta             = "mieu_ta"
	ProductionPlanCustomField_chat_luong_kp       = "chat_luong_kp"
	ProductionPlanCustomField_van_chuyen          = "van_chuyen"
	ProductionPlanCustomField_chi_tiet_van_chuyen = "chi_tiet_van_chuyen"
	ProductionPlanCustomField_ghi_chu_van_chuyen  = "ghi_chu_van_chuyen"
	ProductionPlanCustomField_ten_mau_sp          = "ten_mau_sp"
	ProductionPlanCustomField_hinh_mau_sp         = "hinh_mau_sp"
	ProductionPlanCustomField_hinh_thuc_in        = "hinh_thuc_in"
	ProductionPlanCustomField_loai_in             = "loai_in"
	ProductionPlanCustomField_film                = "film"
	ProductionPlanCustomField_hinh_dang           = "hinh_dang"
	ProductionPlanCustomField_keo_dan             = "keo_dan"
	ProductionPlanCustomField_ghi_chu_keo_dan     = "ghi_chu_keo_dan"
	ProductionPlanCustomField_bdc                 = "bdc"
	ProductionPlanCustomField_epoxy               = "epoxy"
	ProductionPlanCustomField_ghi_chu_epoxy       = "ghi_chu_epoxy"
)

type Service interface {
	FindProductionPlans(ctx context.Context, opts *FindProductionPlansOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error)
	FindProductionPlansWithNoPermission(ctx context.Context, opts *FindProductionPlansOpts, sort *repository.Sort, limit, offset int64) ([]*DataWithNoPermission, *repository.CountResult, error)
	CreateProductionPlan(ctx context.Context, opt *CreateProductionPlanOpts) (string, error)
	EditProductionPlan(ctx context.Context, opt *EditProductionPlanOpts) error
	DeleteProductionPlan(ctx context.Context, id string) error
	ProcessProductionOrder(ctx context.Context, opt *ProcessProductionOrderOpts) (string, error)
}

type productionPlanService struct {
	productionPlanRepo       repository.ProductionPlanRepo
	productionOrderRepo      repository.ProductionOrderRepo
	productionOrderStageRepo repository.ProductionOrderStageRepo
	customFieldRepo          repository.CustomFieldRepo
	customerRepo             repository.CustomerRepo
	userRepo                 repository2.UserRepo
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
		roleService:              roleService,
		cfg:                      cfg,
		redisDB:                  redisDB,
	}
}

type Data struct {
	*repository.ProductionPlanData
	CustomData   map[string]string
	CustomerData *repository.CustomerData
}

type DataWithNoPermission struct {
	*repository.ProductionPlanData
	CustomerData *repository.CustomerData
}
