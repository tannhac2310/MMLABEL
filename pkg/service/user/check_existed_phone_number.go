package user

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

func (u *userService) CheckExistedUserName(ctx context.Context, userName string) (bool, error) {
	user, err := u.userNamePasswordRepo.FindByUserName(ctx, userName)
	if err != nil && !cockroach.IsErrNoRows(err) {
		return false, fmt.Errorf("u.userNamePasswordRepo.FindByUserName: %w", err)
	}

	return user != nil, nil
}
