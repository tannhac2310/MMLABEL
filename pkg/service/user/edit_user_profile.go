package user

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type EditUserProfileOpts struct {
	ID          string
	Name        string
	PhoneNumber string
	Avatar      string
	Status      enum.UserStatus
	Address     string
	Email       string
}

func (u *userService) EditUserProfile(ctx context.Context, opts *EditUserProfileOpts) error {
	var err error
	table := model.User{}
	updater := cockroach.NewUpdater(table.TableName(), model.UserFieldID, opts.ID)

	updater.Set(model.UserFieldName, opts.Name)
	updater.Set(model.UserFieldStatus, opts.Status)
	updater.Set(model.UserFieldPhoneNumber, opts.PhoneNumber)
	updater.Set(model.UserFieldAvatar, opts.Avatar)
	updater.Set(model.UserFieldAddress, opts.Address)
	updater.Set(model.UserFieldEmail, opts.Email)

	err = cockroach.UpdateFields(ctx, updater)
	if err != nil {
		return fmt.Errorf("err update user status: %w", err)
	}

	return nil

}
