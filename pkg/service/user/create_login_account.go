package user

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/passwordutil"
)

func (u *userService) CreateLoginAccount(ctx context.Context, userID, email, phoneNumber, password string) error {
	now := time.Now()
	pwd, _ := passwordutil.HashAndSalt(password)

	userNamePassword := &model.UserNamePassword{
		ID:          idutil.ULIDNow(),
		UserID:      userID,
		Email:       cockroach.String(email),
		PhoneNumber: cockroach.String(phoneNumber),
		Password:    pwd,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	err := u.userNamePasswordRepo.Insert(ctx, userNamePassword)
	if err != nil {
		return fmt.Errorf("userNamePasswordRepo.Insert: %w", err)
	}

	return nil
}
