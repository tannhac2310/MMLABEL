package auth

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/user"
)

type LoginFirebaseResult struct {
	UserID       string
	Token        string
	RefreshToken string
	ACL          []string
}

func (a *authService) LoginFirebase(ctx context.Context, idToken string) (*LoginResult, error) {
	token, err := a.firebaseAuth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("error verifying ID token: %v", err)
	}

	userFirebase, err := a.userFirebaseRepo.FindByID(ctx, token.Subject)
	if err != nil && !cockroach.IsErrNoRows(err) {
		return nil, fmt.Errorf("a.userFirebaseRepo.FindByID: %w", err)
	}

	userID := ""
	// storage user in our system
	if userFirebase != nil {
		userID = userFirebase.UserID
	} else {
		err := cockroach.ExecInTx(ctx, func(ctx context.Context) error {
			u, err := a.firebaseAuth.GetUser(ctx, token.Subject)
			if err != nil {
				return fmt.Errorf("a.firebaseAuth.GetUser: %w", err)
			}

			userID, err = a.userService.CreateFirebaseUser(ctx, token.Subject, &user.CreateUserOpts{
				Name:        u.DisplayName,
				Avatar:      u.PhotoURL,
				PhoneNumber: u.PhoneNumber,
				Email:       u.Email,
			})
			if err != nil {
				return fmt.Errorf("a.userService.CreateFirebaseUser: %w", err)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("cockroach.ExecInTx: %w", err)
		}
	}

	return a.buildLoginResult(ctx, userID)
}
