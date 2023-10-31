package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/dto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

func (u *userController) FindUserByID(c *gin.Context) {

	req := &dto.FindUserByIDRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	user, err := u.userService.FindUserByID(c, req.ID)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	groupIDs, err := u.groupService.GetGroupsForUser(c, user.ID)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	roleIDs, err := u.roleService.GetRolesForUser(c, user.ID)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.FindUserByIDResponse{
		ID:          user.ID,
		Name:        user.Name,
		Avatar:      user.Avatar,
		Address:     user.Address,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Type:        user.Type,
		Status:      user.Status,
		GroupIDs:    groupIDs,
		RoleIDs:     roleIDs,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	})
}
