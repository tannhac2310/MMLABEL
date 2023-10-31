package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/dto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/user"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

func (u *userController) CreateUser(c *gin.Context) {
	req := &dto.CreateUserRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	userID, err := u.userService.CreateUser(c, &user.CreateUserOpts{
		Name:        req.Name,
		Avatar:      req.Avatar,
		Address:     req.Address,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Type:        req.Type,
		Roles:       []string{model.UserRoleUser},
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	err = u.userService.CreateLoginAccount(c, userID, req.Email, req.PhoneNumber, req.Password)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.CreateUserResponse{
		ID: userID,
	})
}
