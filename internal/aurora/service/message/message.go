package message

import (
	"context"
	"github.com/go-redis/redis"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/configs"
	repository2 "mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/messagerelation"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

type Service interface {
	CreateMessage(ctx context.Context, opt *CreateMessageOpts) (string, error)
	EditMessage(ctx context.Context, opt *EditMessageOpts) error
	FindMessages(ctx context.Context, opts *FindMessagesOpts, sort *repository2.Sort, limit, offset int64) ([]*Data, *repository2.CountResult, error)
	Delete(ctx context.Context, id, userID string) error
	FindOne(ctx context.Context, opts *FindMessagesOpts) (*repository2.MessageData, error)
}

type messageService struct {
	messageRepo            repository2.MessageRepo
	permissionRepo         repository.PermissionRepo
	messageRelationService messagerelation.Service
	cfg                    *configs.Config
	redisDB                redis.Cmdable
}

func NewService(
	messageRepo repository2.MessageRepo,
	cfg *configs.Config,
	permissionRepo repository.PermissionRepo,
	messageRelationService messagerelation.Service,
	redisDB redis.Cmdable,
) Service {
	return &messageService{
		messageRepo:            messageRepo,
		permissionRepo:         permissionRepo,
		messageRelationService: messageRelationService,
		cfg:                    cfg,
		redisDB:                redisDB,
	}
}

type Data struct {
	*repository2.MessageData
}
