package auth

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/user"
)

func (a *authService) LoginOTP(ctx context.Context, name, phoneNumber, passcode string) (*LoginResult, error) {

	u, err := a.userRepo.FindByPhoneOrEmail(ctx, phoneNumber)
	if err != nil && !cockroach.IsErrNoRows(err) {
		return nil, fmt.Errorf("u.userRepo.FindByPhoneOrEmail: %w", err)
	}

	userID := ""
	if u == nil {
		userID, err = a.userService.CreateUser(ctx, &user.CreateUserOpts{
			Name:        name,
			PhoneNumber: phoneNumber,
			Email:       phoneNumber,
		})
		if err != nil {
			return nil, fmt.Errorf("a.userService.CreateUser: %w", err)
		}
	}

	if userID == "" {
		userID = u.ID
	}

	return a.buildLoginResult(ctx, userID)
}
