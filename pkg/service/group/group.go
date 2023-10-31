package group

import (
	"context"
	"github.com/casbin/casbin/v2"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

type Service interface {
	CreateGroup(ctx context.Context, opt *CreateGroupOpts) error
	EditGroup(ctx context.Context, opt *EditGroupOpts) error
	FindGroups(ctx context.Context, opts *FindGroupsOpts, limit, offset int64) ([]*Group, error)
	FindGroupByID(ctx context.Context, id string) (*Group, error)
	AddGroupsForUser(ctx context.Context, userID string, groupIDs []string) error
	RemoveGroupsForUser(ctx context.Context, userID string, groupIDs []string) error
	GetGroupsForUser(ctx context.Context, userID string) ([]string, error)
	GetUsersForGroup(ctx context.Context, groupID string) ([]string, error)
}

type groupService struct {
	groupRepo     repository.GroupRepo
	endforcer     casbin.IEnforcer
	userGroupRepo repository.UserGroupRepo
}

func NewService(
	groupRepo repository.GroupRepo,
	endforcer casbin.IEnforcer,
	userGroupRepo repository.UserGroupRepo,
) Service {
	return &groupService{
		groupRepo:     groupRepo,
		endforcer:     endforcer,
		userGroupRepo: userGroupRepo,
	}
}

type CreateGroupOpts struct {
	ID    string
	Name  string
	Roles []string
}

type EditGroupOpts CreateGroupOpts

type Group struct {
	*model.Group
	Roles []string
}
