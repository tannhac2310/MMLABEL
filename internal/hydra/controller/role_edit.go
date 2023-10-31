package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/dto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/role"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

func (r *roleController) EditRole(c *gin.Context) {

	req := &dto.EditRoleRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = r.roleService.EditRole(c, &role.EditRoleOpts{
		ID:          req.ID,
		Name:        req.Name,
		Priority:    req.Priority,
		Permissions: req.Permissions,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.EditRoleResponse{})
}
