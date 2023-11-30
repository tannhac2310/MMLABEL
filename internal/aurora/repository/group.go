package repository

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type GroupRepo interface {
	Insert(ctx context.Context, e *model.Group) error
	Update(ctx context.Context, e *model.Group) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchGroupOpts) ([]*GroupData, error)
	Count(ctx context.Context, s *SearchGroupOpts) (*CountResult, error)
}

type sGroupRepo struct {
}

func NewGroupRepo() GroupRepo {
	return &sGroupRepo{}
}

func (r *sGroupRepo) Insert(ctx context.Context, e *model.Group) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}
