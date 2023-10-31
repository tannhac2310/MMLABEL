package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/dto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

func (a *authController) RefreshToken(c *gin.Context) {

	req := &dto.RefreshTokenRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	loginResult, err := a.authService.RefreshToken(c, req.RefreshToken)
	if err != nil {
		transportutil.Error(c, apperror.ErrUnauthenticated.WithDebugMessage(err.Error()))
		return
	}

	resp, err := a.toLoginResponse(c, loginResult)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, resp)
}
