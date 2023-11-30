package user

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

type SearchUsersOpts struct {
	IDs         []string
	NotIDs      []string
	NotRoleIDs  []string
	Name        string
	Department	string
	Search      string
	PhoneNumber string
	Email       string
	GroupID     string
	RoleID      string
	Type        enum.UserType
}

func (u *userService) SearchUsers(ctx context.Context, opts *SearchUsersOpts, limit, offset int64) ([]*repository.UserData, *repository.CountResult, error) {
	filter := &repository.SearchUsersOpts{
		IDs:         opts.IDs,
		NotIDs:      opts.NotIDs,
		NotRoleIDs:  opts.NotRoleIDs,
		Department:  opts.Department,
		Name:        opts.Name,
		Search:      opts.Search,
		PhoneNumber: opts.PhoneNumber,
		Email:       opts.Email,
		GroupID:     opts.GroupID,
		RoleID:      opts.RoleID,
		Type:        opts.Type,
		Limit:       limit,
		Offset:      offset,
	}
	results, err := u.userRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	total, err := u.userRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	return results, total, nil
}
