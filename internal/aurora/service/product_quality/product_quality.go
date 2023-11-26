package product_quality

import (
	"context"

	"github.com/go-redis/redis"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/configs"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

type Service interface {
	CreateProductQuality(ctx context.Context, opt *CreateProductQualityOpts) (string, error)
	EditProductQuality(ctx context.Context, opt *EditProductQualityOpts) error
	FindProductQuality(ctx context.Context, opts *FindProductQualityOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, []*ProductQualityAnalysis, error)
	Delete(ctx context.Context, id string) error
}
type ProductQualityAnalysis struct {
	DefectType string `json:"defectType"`
	Count      int64  `json:"count"`
}
type productQualityService struct {
	productQualityRepo repository.ProductQualityRepo
	cfg                *configs.Config
	redisDB            redis.Cmdable
}

func NewService(
	productQualityRepo repository.ProductQualityRepo,
	cfg *configs.Config,
	redisDB redis.Cmdable,
) Service {
	return &productQualityService{
		productQualityRepo: productQualityRepo,
		cfg:                cfg,
		redisDB:            redisDB,
	}
}

type Data struct {
	*repository.ProductQualityData
}
