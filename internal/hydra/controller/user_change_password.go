package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/dto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/interceptor"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

func (u *userController) ChangePassword(c *gin.Context) {

	req := &dto.ChangePasswordRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	userID := interceptor.UserIDFromCtx(c)
	verified, err := u.userService.VerifyPassword(c, userID, req.OldPassword)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	if !verified {
		transportutil.Error(c, apperror.ErrInvalidOldPassword)
		return
	}

	err = u.userService.ChangePassword(c, userID, req.Password)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.ChangePasswordResponse{})
}
