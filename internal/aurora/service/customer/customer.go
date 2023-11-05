package customer

import (
	"context"
	"github.com/go-redis/redis"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/configs"
	repository2 "mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

type Service interface {
	CreateCustomer(ctx context.Context, opt *CreateCustomerOpts) (string, error)
	EditCustomer(ctx context.Context, opt *EditCustomerOpts) error
	FindCustomers(ctx context.Context, opts *FindCustomersOpts, sort *repository2.Sort, limit, offset int64) ([]*Data, *repository2.CountResult, error)
	Delete(ctx context.Context, id, userID string) error
}

type customerService struct {
	customerRepo   repository2.CustomerRepo
	permissionRepo repository.PermissionRepo
	cfg            *configs.Config
	redisDB        redis.Cmdable
}

func NewService(
	customerRepo repository2.CustomerRepo,
	cfg *configs.Config,
	permissionRepo repository.PermissionRepo,
	redisDB redis.Cmdable,
) Service {
	return &customerService{
		customerRepo:   customerRepo,
		permissionRepo: permissionRepo,
		cfg:            cfg,
		redisDB:        redisDB,
	}
}

type Data struct {
	*repository2.CustomerData
}
