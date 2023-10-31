package user

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
)

func (u *userService) FindUserByID(ctx context.Context, id string) (*model.User, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.ErrUserNotFound.WithDebugMessage(err.Error())
	}

	return user, nil
}
