package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/dto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/role"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

func (r *roleController) FindRole(c *gin.Context) {
	req := &dto.FindRoleRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	roles, err := r.roleService.FindRoles(
		c,
		&role.FindRolesOpts{
			IDs:  req.Filter.IDs,
			Name: req.Filter.Name,
		},
		req.Paging.Limit+1,
		req.Paging.Offset)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	roleResp := make([]*dto.Role, 0, len(roles))
	for _, g := range roles {
		roleResp = append(roleResp, &dto.Role{
			ID:          g.ID,
			Name:        g.Name,
			Priority:    g.Priority,
			Permissions: g.Permissions,
			UserCount:   g.UserCount,
		})
	}

	nextPage := &commondto.Paging{
		Limit:  req.Paging.Limit,
		Offset: req.Paging.Offset + req.Paging.Limit,
	}

	if int64(len(roles)) <= req.Paging.Limit {
		nextPage = nil
	}

	if l := int64(len(roleResp)); l > req.Paging.Limit {
		roleResp = roleResp[:req.Paging.Limit]
	}

	transportutil.SendJSONResponse(c, &dto.FindRoleResponse{
		Roles:    roleResp,
		NextPage: nextPage,
	})
}
