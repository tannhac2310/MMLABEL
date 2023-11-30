package repository

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type EventLogRepo interface {
	Insert(ctx context.Context, e *model.EventLog) error
	Update(ctx context.Context, e *model.EventLog) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchEventLogOpts) ([]*EventLogData, error)
	Count(ctx context.Context, s *SearchEventLogOpts) (*CountResult, error)
}

type sEventLogRepo struct {
}

func NewEventLogRepo() EventLogRepo {
	return &sEventLogRepo{}
}

func (r *sEventLogRepo) Insert(ctx context.Context, e *model.EventLog) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}
