package role

import (
	"context"
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"time"

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
	UpsertRolePermissions(ctx context.Context, roleID string, permissions []*Permission) error
	FindRolePermissions(ctx context.Context, roleId string) ([]*repository.RolePermissionData, error)
	FindRolePermissionsByUser(ctx context.Context, userID string) ([]*repository.RolePermissionData, error)
}
type Permission struct {
	EntityType string
	EntityID   string
}
type roleService struct {
	endforcer casbin.IEnforcer

	roleRepo       repository.RoleRepo
	userRoleRepo   repository.UserRoleRepo
	rolePermission repository.RolePermissionRepo
}

func (r *roleService) DeleteRole(ctx context.Context, id string) error {
	err := cockroach.ExecInTx(ctx, func(ctx context.Context) error {
		err := r.roleRepo.Delete(ctx, id)
		if err != nil {
			return fmt.Errorf("DeleteRole r.roleRepo.Delete: %w", err)
		}
		err = r.userRoleRepo.DeleteByRoleID(ctx, id)
		if err != nil {
			return fmt.Errorf("DeleteRole r.userRoleRepo.DeleteByRoleID: %w", err)
		}
		err = r.rolePermission.DeleteByRoleID(ctx, id)
		if err != nil {
			return fmt.Errorf("DeleteRole r.rolePermission.DeleteByRoleID: %w", err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("DeleteRole cockroach.ExecInTx: %w", err)
	}

	return nil
}

func NewService(
	endforcer casbin.IEnforcer,
	roleRepo repository.RoleRepo,
	userRoleRepo repository.UserRoleRepo,
	rolePermission repository.RolePermissionRepo,
) Service {
	return &roleService{
		endforcer:      endforcer,
		roleRepo:       roleRepo,
		rolePermission: rolePermission,
		userRoleRepo:   userRoleRepo,
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

func (s *roleService) FindRolePermissions(ctx context.Context, roleId string) ([]*repository.RolePermissionData, error) {
	return s.rolePermission.Search(ctx, &repository.SearchRolePermissionOpts{
		RoleID: roleId,
		Limit:  10000,
		Offset: 0,
	})
}
func (s *roleService) FindRolePermissionsByUser(ctx context.Context, userId string) ([]*repository.RolePermissionData, error) {
	filter := &repository.SearchUserRoleOpts{
		UserIDs:   []string{userId, "-1"},
		UserTypes: nil,
		Search:    "",
		Limit:     1000,
		Offset:    0,
	}
	userRoles, err := s.userRoleRepo.Search(ctx, filter)

	if err != nil {
		return nil, err
	}

	permissions := make([]*repository.RolePermissionData, 0)

	for _, userRole := range userRoles {
		rolePermissions, err := s.rolePermission.Search(ctx, &repository.SearchRolePermissionOpts{
			RoleID: userRole.RoleID,
			Limit:  10000,
			Offset: 0,
		})
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, rolePermissions...)
	}

	return permissions, nil
}
func (s *roleService) UpsertRolePermissions(ctx context.Context, roleID string, permissions []*Permission) error {

	// delete all permissions for role
	err := cockroach.ExecInTx(ctx, func(ctx context.Context) error {
		err := s.rolePermission.DeleteByRoleID(ctx, roleID)
		if err != nil && !errors.Is(err, repository.ErrNotFound) {
			return err
		}
		for _, p := range permissions {
			err = s.rolePermission.Insert(ctx, &model.RolePermission{
				RoleID:     roleID,
				EntityType: p.EntityType,
				EntityID:   p.EntityID,
				CreatedAt:  time.Now(),
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("UpsertRolePermissions cockroach.ExecInTx: %w", err)
	}
	return nil
}
