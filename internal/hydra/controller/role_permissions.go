package controller

import (
	"github.com/gin-gonic/gin"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/role"

	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/dto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

func (r *roleController) FindRolePermissions(c *gin.Context) {
	req := &dto.FindRolePermissionsRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}
	permissions, err := r.roleService.FindRolePermissions(c.Request.Context(), req.RoleID)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	rolePermission := make([]*dto.RolePermission, len(permissions))
	for i, v := range permissions {
		rolePermission[i] = &dto.RolePermission{
			ID:         v.ID,
			RoleID:     v.RoleID,
			EntityType: v.EntityType,
			EntityID:   v.EntityID,
		}
	}

	transportutil.SendJSONResponse(c, &dto.FindRolePermissionsResponse{
		RolePermissions: rolePermission,
	})
}
func (r *roleController) UpsertRolePermissions(c *gin.Context) {

	req := &dto.UpsertRolePermissionsRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}
	permission := make([]*role.Permission, len(req.Permissions))
	for i, v := range req.Permissions {
		permission[i] = &role.Permission{
			EntityType: v.EntityType,
			EntityID:   v.EntityID,
		}
	}
	err = r.roleService.UpsertRolePermissions(c.Request.Context(), req.RoleID, permission)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.UpsertRolePermissionsResponse{})
}
