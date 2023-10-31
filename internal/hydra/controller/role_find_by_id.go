package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/dto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

func (r *roleController) FindRoleByID(c *gin.Context) {

	req := &dto.FindRoleByIDRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	role, err := r.roleService.FindRoleByID(c, req.ID)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.FindRoleByIDResponse{
		ID:          role.ID,
		Name:        role.Name,
		Priority:    role.Priority,
		Permissions: role.Permissions,
	})
}
