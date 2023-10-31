package repository

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

type UserFirebaseRepo interface {
	Insert(ctx context.Context, e *model.UserFirebase) error
	FindByID(ctx context.Context, id string) (*model.UserFirebase, error)
}

type userFirebaseRepo struct {
}

func NewUserFirebaseRepo() UserFirebaseRepo {
	return &userFirebaseRepo{}
}

func (r *userFirebaseRepo) Insert(ctx context.Context, e *model.UserFirebase) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("cockroach.Create: %w", err)
	}

	return nil
}

func (r *userFirebaseRepo) FindByID(ctx context.Context, id string) (*model.UserFirebase, error) {
	e := &model.UserFirebase{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("cockroach.FindOne: %w", err)
	}

	return e, nil
}
