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

	transportutil.SendJSONResponse(c, &dto.GetCurrentProfileResponse{
		ID:          user.ID,
		Name:        user.Name,
		Avatar:      user.Avatar,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
	})
}
