package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/dto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

func (a *authController) LoginOTP(c *gin.Context) {

	req := &dto.LoginOTPRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	loginResult, err := a.authService.LoginOTP(c, req.Name, req.PhoneNumber, req.OTP)
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
