package device

import (
	"context"

	"github.com/go-redis/redis"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/configs"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

type Service interface {
	CreateDevice(ctx context.Context, opt *CreateDeviceOpts) (string, error)
	EditDevice(ctx context.Context, opt *EditDeviceOpts) error
	FindDevices(ctx context.Context, opts *FindDevicesOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error)
	Delete(ctx context.Context, id string) error
}

type deviceService struct {
	deviceRepo repository.DeviceRepo
	cfg        *configs.Config
	redisDB    redis.Cmdable
}

func NewService(
	deviceRepo repository.DeviceRepo,
	cfg *configs.Config,
	redisDB redis.Cmdable,
) Service {
	return &deviceService{
		deviceRepo: deviceRepo,
		cfg:        cfg,
		redisDB:    redisDB,
	}
}

type Data struct {
	*repository.DeviceData
}
