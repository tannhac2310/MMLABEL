package user

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/passwordutil"
)

func (u *userService) ChangePassword(ctx context.Context, userID, password string) error {
	usernamePwd, err := u.userNamePasswordRepo.FindByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("u.userNamePasswordRepo.FindByUserID: %w", err)
	}

	usernamePwd.Password, err = passwordutil.HashAndSalt(password)
	if err != nil {
		return fmt.Errorf("passwordutil.HashAndSalt: %w", err)
	}

	if err := u.userNamePasswordRepo.Update(ctx, usernamePwd); err != nil {
		return fmt.Errorf("u.usernamePwd.Update: %w", err)
	}

	return nil
}
