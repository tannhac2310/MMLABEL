package department

import (
	"context"

	"github.com/go-redis/redis"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/configs"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

type Service interface {
	CreateDepartment(ctx context.Context, opt *CreateDepartmentOpts) (string, error)
	EditDepartment(ctx context.Context, opt *EditDepartmentOpts) error
	FindDepartments(ctx context.Context, opts *FindDepartmentsOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error)
	Delete(ctx context.Context, id string) error
}

type departmentService struct {
	departmentRepo repository.DepartmentRepo
	cfg            *configs.Config
	redisDB        redis.Cmdable
}

func NewService(
	departmentRepo repository.DepartmentRepo,
	cfg *configs.Config,
	redisDB redis.Cmdable,
) Service {
	return &departmentService{
		departmentRepo: departmentRepo,
		cfg:            cfg,
		redisDB:        redisDB,
	}
}

type Data struct {
	*repository.DepartmentData
}
