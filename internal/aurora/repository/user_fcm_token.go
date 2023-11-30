package repository

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type UserFcmTokenRepo interface {
	Insert(ctx context.Context, e *model.UserFcmToken) error
	Update(ctx context.Context, e *model.UserFcmToken) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchUserFcmTokenOpts) ([]*UserFcmTokenData, error)
	Count(ctx context.Context, s *SearchUserFcmTokenOpts) (*CountResult, error)
}

type sUserFcmTokenRepo struct {
}

func NewUserFcmTokenRepo() UserFcmTokenRepo {
	return &sUserFcmTokenRepo{}
}

func (r *sUserFcmTokenRepo) Insert(ctx context.Context, e *model.UserFcmToken) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}
