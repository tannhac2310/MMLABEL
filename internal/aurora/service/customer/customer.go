package customer

import (
	"context"

	"github.com/go-redis/redis"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/configs"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

type Service interface {
	CreateCustomer(ctx context.Context, opt *CreateCustomerOpts) (string, error)
	EditCustomer(ctx context.Context, opt *EditCustomerOpts) error
	FindCustomers(ctx context.Context, opts *FindCustomersOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error)
	Delete(ctx context.Context, id string) error
}

type customerService struct {
	customerRepo repository.CustomerRepo
	cfg          *configs.Config
	redisDB      redis.Cmdable
}

func NewService(
	customerRepo repository.CustomerRepo,
	cfg *configs.Config,
	redisDB redis.Cmdable,
) Service {
	return &customerService{
		customerRepo: customerRepo,
		cfg:          cfg,
		redisDB:      redisDB,
	}
}

type Data struct {
	*repository.CustomerData
}
