package production_order

import (
	"context"
	"github.com/go-redis/redis"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"

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
	productionOrderRepo      repository.ProductionOrderRepo
	productionOrderStageRepo repository.ProductionOrderStageRepo
	cfg                      *configs.Config
	redisDB                  redis.Cmdable
}

func (c *productionOrderService) deleteProductionOrderStage(ctx context.Context, ids []string, productionId string) interface{} {
	// find production order stage by production order id
	productionOrderStages, err := c.productionOrderStageRepo.Search(ctx, &repository.SearchProductionOrderStagesOpts{
		ProductionOrderID: productionId,
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
	cfg *configs.Config,
	redisDB redis.Cmdable,
) Service {
	return &productionOrderService{
		productionOrderRepo:      productionOrderRepo,
		productionOrderStageRepo: productionOrderStageRepo,
		cfg:                      cfg,
		redisDB:                  redisDB,
	}
}

type Data struct {
	*repository.ProductionOrderData
	ProductionOrderStage []*model.ProductionOrderStage
}
