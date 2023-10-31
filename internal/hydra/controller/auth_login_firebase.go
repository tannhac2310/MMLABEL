package controller

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/dto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/auth"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

func (a *authController) LoginFirebase(c *gin.Context) {

	req := &dto.LoginFirebaseRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	loginResult, err := a.authService.LoginFirebase(c, req.IDToken)
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

func (a *authController) toLoginResponse(c context.Context, r *auth.LoginResult) (*dto.LoginResponse, error) {
	u, err := a.userService.FindUserByID(c, r.UserID)
	fmt.Println("------------------->", u)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:        r.Token,
		RefreshToken: r.RefreshToken,
		Profile: &dto.UserProfile{
			ID:          u.ID,
			Name:        u.Name,
			Avatar:      u.Avatar,
			PhoneNumber: u.PhoneNumber,
			Email:       u.Email,
		},
		ACL:  r.ACL,
		Role: r.MainRole.Name,
	}, nil
}
