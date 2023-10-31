package repository

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

type UserNamePasswordRepo interface {
	Insert(ctx context.Context, e *model.UserNamePassword) error
	FindByID(ctx context.Context, id string) (*model.UserNamePassword, error)
	FindByUserID(ctx context.Context, id string) (*model.UserNamePassword, error)
	FindByUserName(ctx context.Context, userName string) (*model.UserNamePassword, error)
	Update(ctx context.Context, e *model.UserNamePassword) error
}

type userNamePasswordRepo struct {
}

func NewUserNamePasswordRepo() UserNamePasswordRepo {
	return &userNamePasswordRepo{}
}

func (u *userNamePasswordRepo) Insert(ctx context.Context, e *model.UserNamePassword) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("cockroach.Create: %w", err)
	}

	return nil
}

func (u *userNamePasswordRepo) FindByID(ctx context.Context, id string) (*model.UserNamePassword, error) {
	e := &model.UserNamePassword{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("cockroach.FindOne: %w", err)
	}

	return e, nil
}

func (u *userNamePasswordRepo) FindByUserID(ctx context.Context, id string) (*model.UserNamePassword, error) {
	e := &model.UserNamePassword{}
	err := cockroach.FindOne(ctx, e, "user_id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("cockroach.FindOne: %w", err)
	}

	return e, nil
}

func (u *userNamePasswordRepo) FindByUserName(ctx context.Context, userName string) (*model.UserNamePassword, error) {
	e := &model.UserNamePassword{}
	err := cockroach.FindOne(ctx, e, "email = $1 OR phone_number = $1", userName)
	if err != nil {
		return nil, fmt.Errorf("cockroach.FindOne: %w", err)
	}

	return e, nil
}

func (u *userNamePasswordRepo) Update(ctx context.Context, e *model.UserNamePassword) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}
