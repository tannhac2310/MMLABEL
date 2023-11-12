package stage

import (
	"context"

	"github.com/go-redis/redis"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/configs"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

type Service interface {
	CreateStage(ctx context.Context, opt *CreateStageOpts) (string, error)
	EditStage(ctx context.Context, opt *EditStageOpts) error
	FindStages(ctx context.Context, opts *FindStagesOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error)
	Delete(ctx context.Context, id string) error
}

type stageService struct {
	stageRepo repository.StageRepo
	cfg       *configs.Config
	redisDB   redis.Cmdable
}

func NewService(
	stageRepo repository.StageRepo,
	cfg *configs.Config,
	redisDB redis.Cmdable,
) Service {
	return &stageService{
		stageRepo: stageRepo,
		cfg:       cfg,
		redisDB:   redisDB,
	}
}

type Data struct {
	*repository.StageData
}
