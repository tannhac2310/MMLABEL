package auth

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/passwordutil"
)

type LoginUserNamePasswordResult struct {
	UserID       string
	Token        string
	RefreshToken string
	ACL          []string
}

func (a *authService) LoginUserNamePassword(ctx context.Context, userName, password string) (*LoginResult, error) {
	userNamePassword, err := a.userNamePasswordRepo.FindByUserName(ctx, userName)
	if err != nil {
		if cockroach.IsErrNoRows(err) {
			return nil, ErrNotFoundUsername
		}

		return nil, fmt.Errorf("a.userNamePasswordRepo.FindByUserName: %w", err)
	}

	if !passwordutil.ComparePasswords(userNamePassword.Password, password) {
		return nil, ErrWrongPassword
	}

	return a.buildLoginResult(ctx, userNamePassword.UserID)
}
