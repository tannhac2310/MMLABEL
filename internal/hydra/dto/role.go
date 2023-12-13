package dto

import (
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type AddRoleToUsersRequest struct {
	UserIDs []string `json:"userIds" binding:"required"`
	RoleID  string   `json:"roleId" binding:"required"`
}
type RemoveRoleToUsersRequest struct {
	UserIDs []string `json:"userIds" binding:"required"`
	RoleID  string   `json:"roleId" binding:"required"`
}
type Role struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Priority    int      `json:"priority"`
	Permissions []string `json:"permissions"`
	UserCount   int64    `json:"userCount"`
}
type FindPermisionsRequest struct {
}

type FindPermissionsResponse struct {
	HydraPermissions []string `json:"hydraPermissions"`
	GezuPermissions  []string `json:"gezuPermissions"`
}
type DeleteRoleRequest struct {
	ID string `json:"id" binding:"required"`
}
type DeleteRoleResponse struct {
}
type CreateRoleRequest struct {
	ID          string   `json:"id" binding:"required"`
	Name        string   `json:"name" binding:"required"`
	Priority    int      `json:"priority" binding:"required"`
	Permissions []string `json:"permissions"`
}
type CreateRoleResponse struct {
	ID string `json:"id"`
}
type EditRoleRequest struct {
	ID          string   `json:"id" binding:"required"`
	Name        string   `json:"name" binding:"required"`
	Priority    int      `json:"priority" binding:"required"`
	Permissions []string `json:"permissions"`
}
type EditRoleResponse struct {
}
type RoleFilter struct {
	IDs  []string `json:"ids"`
	Name string   `json:"name"`
}
type FindRoleRequest struct {
	Filter *RoleFilter       `json:"filter"`
	Paging *commondto.Paging `json:"paging" binding:"required"`
}
type FindRoleResponse struct {
	Roles    []*Role           `json:"roles"`
	NextPage *commondto.Paging `json:"nextPage"`
}
type FindRoleByIDRequest struct {
	ID string `json:"id" binding:"required"`
}
type FindRoleByIDResponse Role
type AddRolesForUserResponse struct{}
type RemoveRolesForUserResponse struct{}
type RoleUsersFilter struct {
	RoleIDs []string `json:"roleIds" binding:"required"`
	Search  string   `json:"search"`
}
type FindRoleUsersRequest struct {
	Filter *RoleUsersFilter  `json:"filter"`
	Paging *commondto.Paging `json:"paging" binding:"required"`
}

type RoleUsers struct {
	ID            string          `json:"id"`
	RoleID        string          `json:"roleId"`
	RoleName      string          `json:"roleName"`
	Name          string          `json:"name"`
	Avatar        string          `json:"avatar"`
	PhoneNumber   string          `json:"phoneNumber"`
	Email         string          `json:"email"`
	Address       string          `json:"address"`
	Type          enum.UserType   `json:"type"`
	Status        enum.UserStatus `json:"status"`
	CreatedBy     string          `json:"createdBy"`
	CreatedByName string          `json:"createdByName"`
	CreatedAt     time.Time       `json:"createdAt"`
}
type FindRoleUsersResponse struct {
	RoleUsers []*RoleUsers      `json:"roleUsers"`
	NextPage  *commondto.Paging `json:"nextPage"`
	Total     int64             `json:"total"`
}

type RolePermission struct {
	ID         string `json:"id"`
	RoleID     string `json:"roleId"`
	EntityType string `json:"entityType"`
	EntityID   string `json:"entityId"`
}

type FindRolePermissionsRequest struct {
	RoleID string `json:"roleId" binding:"required"`
}

type FindRolePermissionsResponse struct {
	RolePermissions []*RolePermission `json:"rolePermissions"`
}
type Permission struct {
	EntityID   string `json:"entityId"`
	EntityType string `json:"entityType"`
}
type UpsertRolePermissionsRequest struct {
	RoleID      string        `json:"roleId" binding:"required"`
	Permissions []*Permission `json:"permissions" binding:"required"`
}
type UpsertRolePermissionsResponse struct {
}
