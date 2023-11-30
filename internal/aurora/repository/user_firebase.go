package repository

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type UserFirebaseRepo interface {
	Insert(ctx context.Context, e *model.UserFirebase) error
	Update(ctx context.Context, e *model.UserFirebase) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchUserFirebaseOpts) ([]*UserFirebaseData, error)
	Count(ctx context.Context, s *SearchUserFirebaseOpts) (*CountResult, error)
}

type sUserFirebaseRepo struct {
}

func NewUserFirebaseRepo() UserFirebaseRepo {
	return &sUserFirebaseRepo{}
}

func (r *sUserFirebaseRepo) Insert(ctx context.Context, e *model.UserFirebase) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}
