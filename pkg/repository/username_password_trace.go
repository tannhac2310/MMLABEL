package repository

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/tracingutil"
)

type userNamePasswordRepoTrace struct {
	UserNamePasswordRepo
}

func NewUserNamePasswordRepoTrace() UserNamePasswordRepo {
	return &userNamePasswordRepoTrace{
		UserNamePasswordRepo: NewUserNamePasswordRepo(),
	}
}

func (u *userNamePasswordRepoTrace) Insert(ctx context.Context, e *model.UserNamePassword) error {
	ctx, span := tracingutil.Start(ctx, "userNamePasswordRepo.Insert")
	defer span.End()

	return u.UserNamePasswordRepo.Insert(ctx, e)
}

func (u *userNamePasswordRepoTrace) FindByID(ctx context.Context, id string) (*model.UserNamePassword, error) {
	ctx, span := tracingutil.Start(ctx, "userNamePasswordRepo.FindByID")
	defer span.End()

	return u.UserNamePasswordRepo.FindByID(ctx, id)
}

func (u *userNamePasswordRepoTrace) FindByUserID(ctx context.Context, id string) (*model.UserNamePassword, error) {
	ctx, span := tracingutil.Start(ctx, "userNamePasswordRepo.FindByUserID")
	defer span.End()

	return u.UserNamePasswordRepo.FindByUserID(ctx, id)
}

func (u *userNamePasswordRepoTrace) FindByUserName(ctx context.Context, userName string) (*model.UserNamePassword, error) {
	ctx, span := tracingutil.Start(ctx, "userNamePasswordRepo.FindByUserName")
	defer span.End()

	return u.UserNamePasswordRepo.FindByUserName(ctx, userName)
}

func (u *userNamePasswordRepoTrace) Update(ctx context.Context, e *model.UserNamePassword) error {
	ctx, span := tracingutil.Start(ctx, "userNamePasswordRepo.Update")
	defer span.End()

	return u.UserNamePasswordRepo.Update(ctx, e)
}
