package user

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (u *userService) CreateFirebaseUser(ctx context.Context, firebaseID string, opts *CreateUserOpts) (string, error) {
	userID, err := u.CreateUser(ctx, opts)
	if err != nil {
		return "", fmt.Errorf("u.CreateUser: %w", err)
	}

	now := time.Now()
	userFirebase := &model.UserFirebase{
		ID:        firebaseID,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := u.userFirebaseRepo.Insert(ctx, userFirebase); err != nil {
		return "", fmt.Errorf("a.userFirebaseRepo.Insert: %w", err)
	}

	return userID, nil
}
