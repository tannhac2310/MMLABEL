package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/dto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

func (a *authController) ResetPassword(c *gin.Context) {

	req := &dto.ResetPasswordRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	loginResult, err := a.authService.ResetPassword(c, req.PhoneNumber, req.OTP, req.NewPassword)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	resp, err := a.toLoginResponse(c, loginResult)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, resp)
}
