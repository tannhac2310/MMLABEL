package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/dto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/interceptor"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

func (u *userController) GetCurrentProfile(c *gin.Context) {

	userID := interceptor.UserIDFromCtx(c)

	user, err := u.userService.FindUserByID(c, userID)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}
	// get user role
	roleIDs, err := u.roleService.GetRolesForUser(c, user.ID)
	if err != nil {
		transportutil.Error(c, err)
		return
	}
	// get user permission
	rolePermissions := make([]*dto.RolePermission, 0)
	for _, roleID := range roleIDs {
		permissions, err := u.roleService.FindRolePermissions(c, roleID)
		if err != nil {
			transportutil.Error(c, err)
			return
		}
		for _, permission := range permissions {
			rolePermissions = append(rolePermissions, &dto.RolePermission{
				ID:         permission.ID,
				RoleID:     permission.RoleID,
				EntityType: permission.EntityType,
				EntityID:   permission.EntityID,
			})
		}
	}
	transportutil.SendJSONResponse(c, &dto.GetCurrentProfileResponse{
		ID:          user.ID,
		Name:        user.Name,
		Code:        user.Code,
		Departments: user.Departments.String,
		Avatar:      user.Avatar,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Type:        user.Type,
		RoleIDs:     roleIDs,
		Permissions: rolePermissions,
	})
}
