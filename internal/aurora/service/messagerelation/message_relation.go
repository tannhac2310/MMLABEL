package messagerelation

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

type Service interface {
	UpsertMessageRelation(ctx context.Context, opt *UpsertMessageRelationOpts) error
	CreateMessageRelation(ctx context.Context, opt *UpsertMessageRelationOpts) (string, error)
	EditMessageRelation(ctx context.Context, opt *EditMessageRelationOpts) error
	FindMessageRelations(ctx context.Context, opts *FindMessageRelationsOpts, limit, offset int64) ([]*repository.MessageRelationData, *repository.CountResult, error)
	SoftDelete(ctx context.Context, id string) error
	FindOne(ctx context.Context, opts *FindMessageRelationsOpts) (*repository.MessageRelationData, error)
}

type messageRelationService struct {
	messageRelationRepo repository.MessageRelationRepo
}

func NewService(
	messageRelationRepo repository.MessageRelationRepo,
) Service {
	return &messageRelationService{
		messageRelationRepo: messageRelationRepo,
	}
}
