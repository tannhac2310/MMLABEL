package user

import (
	"context"
)

func (u *userService) SoftDelete(ctx context.Context, id string) error {
	return u.userRepo.SoftDelete(ctx, id)
}
