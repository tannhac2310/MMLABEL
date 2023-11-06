package production_order

import (
	"context"

	"github.com/go-redis/redis"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/configs"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

type Service interface {
	CreateProductionOrder(ctx context.Context, opt *CreateProductionOrderOpts) (string, error)
	EditProductionOrder(ctx context.Context, opt *EditProductionOrderOpts) error
	FindProductionOrders(ctx context.Context, opts *FindProductionOrdersOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error)
	Delete(ctx context.Context, id string) error
}

type productionOrderService struct {
	productionOrderRepo repository.ProductionOrderRepo
	cfg                 *configs.Config
	redisDB             redis.Cmdable
}

func NewService(
	productionOrderRepo repository.ProductionOrderRepo,
	cfg *configs.Config,
	redisDB redis.Cmdable,
) Service {
	return &productionOrderService{
		productionOrderRepo: productionOrderRepo,
		cfg:                 cfg,
		redisDB:             redisDB,
	}
}

type Data struct {
	*repository.ProductionOrderData
}
