package auth

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

func (a *authService) ResetPassword(ctx context.Context, userName, passcode, newPassword string) (*LoginResult, error) {

	userID := ""
	usernamePwd, err := a.userNamePasswordRepo.FindByUserName(ctx, userName)
	if err != nil && !cockroach.IsErrNoRows(err) {
		return nil, fmt.Errorf("u.userNamePasswordRepo.FindByUserName: %w", err)
	}

	if usernamePwd == nil {
		u, err := a.userRepo.FindByPhoneOrEmail(ctx, userName)
		if err != nil {
			return nil, fmt.Errorf("a.userRepo.FindByPhoneOrEmail: %w", err)
		}

		err = a.userService.CreateLoginAccount(ctx, u.ID, u.Email, u.PhoneNumber, newPassword)
		if err != nil {
			return nil, fmt.Errorf("u.CreateLoginAccount: %w", err)
		}

		userID = u.ID
	}

	if userID == "" {
		userID = usernamePwd.UserID
	}

	err = a.userService.ChangePassword(ctx, userID, newPassword)
	if err != nil {
		return nil, fmt.Errorf("u.ChangePassword: %w", err)
	}

	return a.buildLoginResult(ctx, userID)
}
