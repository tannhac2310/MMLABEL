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
	return []string{
		"ten_sp",
		"ma_sp",
		"ma_sp_phu",
		"dvt",
		"dai",
		"rong",
		"don_vi_dvt",
		"so_luong_mau",
		"so_lan_in",
		"chat_lieu_in",
		"ma_chat_lieu_in",
		"vat_lieu_thay_the",
		"keo_2_mat",
		"ma_keo_2_mat",
		"so_luong_keo_2_mat",
		"mieu_ta",
		"chat_luong_kp",
		"van_chuyen",
		"chi_tiet_van_chuyen",
		"ghi_chu_van_chuyen",
		"ten_mau_sp",
		"hinh_mau_sp",
		"loai_hinh",
		"hinh_sp",
		"hinh_thuc_in",
		"loai_in",
		"film",
		"hinh_dang",
		"keo_dan",
		"ghi_chu_keo_dan",
		"bdc",
		"epoxy",
		"ghi_chu_epoxy",
	}
}

type Data struct {
	*repository.ProductionPlanData
	CustomData map[string]string
}

type DataWithNoPermission struct {
	*repository.ProductionPlanData
}
