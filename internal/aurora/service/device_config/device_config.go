package device_config

import (
	"context"

	"github.com/go-redis/redis"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/configs"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

type Service interface {
	CreateDeviceConfig(ctx context.Context, opt *CreateDeviceConfigOpts) (string, error)
	EditDeviceConfig(ctx context.Context, opt *EditDeviceConfigOpts) error
	FindDeviceConfigs(ctx context.Context, opts *FindDeviceConfigsOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error)
	Delete(ctx context.Context, id string) error
}

type deviceConfigService struct {
	deviceConfigRepo repository.ProductionOrderDeviceConfigRepo
	cfg              *configs.Config
	redisDB          redis.Cmdable
}

func NewService(
	deviceConfigRepo repository.ProductionOrderDeviceConfigRepo,
	cfg *configs.Config,
	redisDB redis.Cmdable,
) Service {
	return &deviceConfigService{
		deviceConfigRepo: deviceConfigRepo,
		cfg:              cfg,
		redisDB:          redisDB,
	}
}

type Data struct {
	*repository.ProductionOrderDeviceConfigData
}
