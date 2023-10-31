package role

import (
	"context"
	"github.com/casbin/casbin/v2"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

type Service interface {
	CreateRole(ctx context.Context, opt *CreateRoleOpts) error
	DeleteRole(ctx context.Context, id string) error
	EditRole(ctx context.Context, opt *EditRoleOpts) error
	FindRoles(ctx context.Context, opts *FindRolesOpts, limit, offset int64) ([]*Role, error)
	FindRoleByID(ctx context.Context, id string) (*Role, error)
	AddRolesForUser(ctx context.Context, userID string, roleIDs []string) error
	AddRoleToUsers(ctx context.Context, userIDs []string, roleID, createdBy string) error
	RemoveRolesForUser(ctx context.Context, userID string, roleIDs []string) error
	RemoveUsers(ctx context.Context, userIDs []string, roleID string) error
	GetRolesForUser(ctx context.Context, userID string) ([]string, error)
	GetUsersForRole(ctx context.Context, opts *FindRoleUsersOpts, limit, offset int64) ([]*repository.RuleUsersData, *repository.CountResult, error)
	HighestRole(ctx context.Context, ids []string) (*model.Role, error)
	AddPolicy(params ...string) (bool, error)
	RemovePolicy(params ...string) (bool, error)
}

type roleService struct {
	endforcer casbin.IEnforcer

	roleRepo     repository.RoleRepo
	userRoleRepo repository.UserRoleRepo
}

func (r *roleService) DeleteRole(ctx context.Context, id string) error {
	return r.roleRepo.SoftDelete(ctx, id)
}

func NewService(
	endforcer casbin.IEnforcer,
	roleRepo repository.RoleRepo,
	userRoleRepo repository.UserRoleRepo,
) Service {
	return &roleService{
		endforcer:    endforcer,
		roleRepo:     roleRepo,
		userRoleRepo: userRoleRepo,
	}
}

type CreateRoleOpts struct {
	ID          string
	Name        string
	Priority    int
	Permissions []string
}

type EditRoleOpts CreateRoleOpts

type Role struct {
	*model.Role
	Permissions []string
	UserCount   int64
}
