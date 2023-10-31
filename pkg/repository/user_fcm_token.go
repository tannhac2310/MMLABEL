package repository

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

type UserFCMTokenRepo interface {
	Insert(ctx context.Context, e *model.UserFCMToken) error
	Update(ctx context.Context, e *model.UserFCMToken) error
	FindByID(ctx context.Context, id string) (*model.UserFCMToken, error)
	FindByUserIDAndDeviceID(ctx context.Context, userID, deviceID string) (*model.UserFCMToken, error)
	FindByUserIDs(ctx context.Context, ids []string) ([]*model.UserFCMToken, error)
}

type userFCMTokenRepo struct {
}

func NewUserFCMTokenRepo() UserFCMTokenRepo {
	return &userFCMTokenRepo{}
}

func (r *userFCMTokenRepo) Insert(ctx context.Context, e *model.UserFCMToken) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("cockroach.Create: %w", err)
	}

	return nil
}

func (r *userFCMTokenRepo) FindByID(ctx context.Context, id string) (*model.UserFCMToken, error) {
	e := &model.UserFCMToken{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("cockroach.FindOne: %w", err)
	}

	return e, nil
}

func (r *userFCMTokenRepo) FindByUserIDs(ctx context.Context, ids []string) ([]*model.UserFCMToken, error) {
	result := make([]*model.UserFCMToken, 0)
	err := cockroach.FindMany(ctx, &model.UserFCMToken{}, &result, "user_id = ANY($1)", ids)
	if err != nil {
		return nil, fmt.Errorf("cockroach.FindMany: %w", err)
	}

	return result, nil
}

func (r *userFCMTokenRepo) FindByUserIDAndDeviceID(ctx context.Context, userID, deviceID string) (*model.UserFCMToken, error) {
	e := &model.UserFCMToken{}
	err := cockroach.FindOne(ctx, e, "user_id = $1 AND device_id = $2", userID, deviceID)
	if err != nil {
		return nil, fmt.Errorf("cockroach.FindOne: %w", err)
	}

	return e, nil
}

func (r *userFCMTokenRepo) Update(ctx context.Context, e *model.UserFCMToken) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}
