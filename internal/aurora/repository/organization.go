package repository

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type OrganizationRepo interface {
	Insert(ctx context.Context, e *model.Organization) error
	Update(ctx context.Context, e *model.Organization) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchOrganizationOpts) ([]*OrganizationData, error)
	Count(ctx context.Context, s *SearchOrganizationOpts) (*CountResult, error)
}

type sOrganizationRepo struct {
}

func NewOrganizationRepo() OrganizationRepo {
	return &sOrganizationRepo{}
}

func (r *sOrganizationRepo) Insert(ctx context.Context, e *model.Organization) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}
