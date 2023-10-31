package user

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/passwordutil"
)

func (u *userService) VerifyPassword(ctx context.Context, userID, password string) (bool, error) {
	usernamePwd, err := u.userNamePasswordRepo.FindByUserID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("u.userNamePasswordRepo.FindByUserID: %w", err)
	}

	return passwordutil.ComparePasswords(usernamePwd.Password, password), nil
}
