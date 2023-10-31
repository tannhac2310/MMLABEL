package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/dto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/user"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

func (u *userController) EditUser(c *gin.Context) {
	req := &dto.EditUserRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = u.userService.EditUserProfile(c, &user.EditUserProfileOpts{
		ID:          req.ID,
		Name:        req.Name,
		Avatar:      req.Avatar,
		Status:      req.Status,
		Address:     req.Address,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.EditUserResponse{})
}
