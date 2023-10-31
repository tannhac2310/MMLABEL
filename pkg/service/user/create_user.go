package user

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

type CreateUserOpts struct {
	Name        string
	Avatar      string
	PhoneNumber string
	Email       string
	Address     string
	Type        enum.UserType
	Roles       []string
}

func (u *userService) CreateUser(ctx context.Context, opts *CreateUserOpts) (string, error) {
	now := time.Now()

	user := &model.User{
		ID:          idutil.ULIDNow(),
		Name:        opts.Name,
		Avatar:      opts.Avatar,
		Address:     opts.Address,
		PhoneNumber: opts.PhoneNumber,
		Email:       opts.Email,
		Status:      enum.UserStatusActive,
		Type:        opts.Type,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := u.userRepo.Insert(ctx, user); err != nil {
		return "", fmt.Errorf("a.userRepo.Insert: %w", err)
	}

	err := u.roleService.AddRolesForUser(ctx, user.ID, opts.Roles)
	if err != nil {
		return "", fmt.Errorf("u.roleService.AddUsers: %w", err)
	}

	return user.ID, nil
}
