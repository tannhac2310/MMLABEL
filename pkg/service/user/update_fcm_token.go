package user

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (u *userService) UpdateFCMToken(ctx context.Context, userID, deviceID, token string) error {
	fcmToken, err := u.userFCMTokenRepo.FindByUserIDAndDeviceID(ctx, userID, deviceID)
	if err != nil && !cockroach.IsErrNoRows(err) {
		return fmt.Errorf("u.userFCMTokenRepo.FindByUserIDAndDeviceID: %w", err)
	}

	if fcmToken == nil {
		now := time.Now()
		fcmToken = &model.UserFCMToken{
			ID:        idutil.ULIDNow(),
			UserID:    userID,
			DeviceID:  deviceID,
			Token:     token,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := u.userFCMTokenRepo.Insert(ctx, fcmToken); err != nil {
			return fmt.Errorf("u.userFCMTokenRepo.Insert: %w", err)
		}
	} else {
		fcmToken.Token = token
		if err := u.userFCMTokenRepo.Update(ctx, fcmToken); err != nil {
			return fmt.Errorf("u.userFCMTokenRepo.Update: %w", err)
		}
	}

	return nil
}
