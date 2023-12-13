package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/dto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/constants"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/role"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

type RoleController interface {
	CreateRole(c *gin.Context)
	EditRole(c *gin.Context)
	FindRole(c *gin.Context)
	FindRoleUser(c *gin.Context)
	FindRoleByID(c *gin.Context)
	AddUsers(c *gin.Context)
	RemoveUsers(c *gin.Context)
	DeleteRole(c *gin.Context)
	UpsertRolePermissions(c *gin.Context)
	FindRolePermissions(c *gin.Context)
	FindPermissions(ctx *gin.Context)
}

func (r *roleController) FindPermissions(ctx *gin.Context) {
	transportutil.SendJSONResponse(ctx, &dto.FindPermissionsResponse{
		HydraPermissions: constants.HydraPermissions,
		GezuPermissions:  constants.GezuPermissionsList,
	})
}

type roleController struct {
	roleService role.Service
}

func RegisterRoleController(
	r *gin.RouterGroup,
	roleService role.Service,
) {
	g := r.Group("role")

	var c RoleController = &roleController{
		roleService: roleService,
	}
	routeutil.AddEndpoint(
		g,
		"find-permissions",
		c.FindPermissions,
		&dto.FindPermisionsRequest{},
		&dto.FindPermissionsResponse{},
		"find permissions",
	)
	routeutil.AddEndpoint(
		g,
		"create-role",
		c.CreateRole,
		&dto.CreateRoleRequest{},
		&dto.CreateRoleResponse{},
		"Create role",
	)
	routeutil.AddEndpoint(
		g,
		"delete-role",
		c.DeleteRole,
		&dto.DeleteRoleRequest{},
		&dto.DeleteRoleResponse{},
		"Delete role",
	)

	routeutil.AddEndpoint(
		g,
		"edit-role",
		c.EditRole,
		&dto.EditRoleRequest{},
		&dto.EditRoleResponse{},
		"Edit role",
	)

	routeutil.AddEndpoint(
		g,
		"find-role",
		c.FindRole,
		&dto.FindRoleRequest{},
		&dto.FindRoleResponse{},
		"Find role",
	)

	routeutil.AddEndpoint(
		g,
		"find-users-in-role",
		c.FindRoleUser,
		&dto.FindRoleUsersRequest{},
		&dto.FindRoleUsersResponse{},
		"Lấy danh sách người dùng trong role",
	)

	routeutil.AddEndpoint(
		g,
		"add-users",
		c.AddUsers,
		&dto.AddRoleToUsersRequest{},
		&dto.AddRolesForUserResponse{},
		"Add users to role",
	)

	routeutil.AddEndpoint(
		g,
		"remove-users",
		c.RemoveUsers,
		&dto.RemoveRoleToUsersRequest{},
		&dto.RemoveRolesForUserResponse{},
		"Remove roles for user",
	)

	routeutil.AddEndpoint(
		g,
		"get-role-permission",
		c.FindRolePermissions,
		&dto.FindRolePermissionsRequest{},
		&dto.FindRolePermissionsResponse{},
		"Remove roles for user",
	)
	routeutil.AddEndpoint(
		g,
		"upsert-role-permission",
		c.UpsertRolePermissions,
		&dto.UpsertRolePermissionsRequest{},
		&dto.UpsertRolePermissionsResponse{},
		"Upsert role permission",
	)
}
