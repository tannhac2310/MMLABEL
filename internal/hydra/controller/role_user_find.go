package controller

import (
	"github.com/gin-gonic/gin"
	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/dto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/role"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

func (r *roleController) FindRoleUser(c *gin.Context) {
	req := &dto.FindRoleUsersRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	usersForRole, total, err := r.roleService.GetUsersForRole(
		c,
		&role.FindRoleUsersOpts{
			RoleIDs: req.Filter.RoleIDs,
			Search:  req.Filter.Search,
		},
		req.Paging.Limit+1,
		req.Paging.Offset)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	roleUsers := make([]*dto.RoleUsers, 0, len(usersForRole))
	for _, g := range usersForRole {
		roleUsers = append(roleUsers, &dto.RoleUsers{
			ID:            g.ID,
			RoleID:        g.RoleID,
			RoleName:      g.RoleName,
			Name:          g.Name,
			Avatar:        g.Avatar,
			PhoneNumber:   g.PhoneNumber,
			Email:         g.Email,
			Address:       g.Address,
			Type:          g.Type,
			Status:        g.Status,
			CreatedBy:     g.CreatedBy.String,
			CreatedByName: g.CreatedByName.String,
			CreatedAt:     g.CreatedAt,
		})
	}

	nextPage := &commondto.Paging{
		Limit:  req.Paging.Limit,
		Offset: req.Paging.Offset + req.Paging.Limit,
	}

	if int64(len(usersForRole)) <= req.Paging.Limit {
		nextPage = nil
	}

	if l := int64(len(roleUsers)); l > req.Paging.Limit {
		roleUsers = roleUsers[:req.Paging.Limit]
	}

	transportutil.SendJSONResponse(c, &dto.FindRoleUsersResponse{
		RoleUsers: roleUsers,
		NextPage:  nextPage,
		Total:     total.Count,
	})
}
