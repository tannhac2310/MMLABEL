package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/dto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/interceptor"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

func (u *userController) RegisterLoginAccount(c *gin.Context) {
	req := &dto.RegisterLoginAccountRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	isExisted, err := u.userService.CheckExistedUserName(c, req.UserName)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	if isExisted {
		transportutil.Error(c, apperror.ErrAccountExisted)
		return
	}

	userID := interceptor.UserIDFromCtx(c)
	err = u.userService.CreateLoginAccount(c, userID, req.UserName, req.UserName, req.Password)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.RegisterLoginAccountResponse{})
}
