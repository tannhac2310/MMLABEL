package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/dto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/group"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

func (g *groupController) FindGroup(c *gin.Context) {

	req := &dto.FindGroupRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	groups, err := g.groupService.FindGroups(
		c,
		&group.FindGroupsOpts{
			IDs:  req.Filter.IDs,
			Name: req.Filter.Name,
		},
		req.Paging.Limit+1,
		req.Paging.Offset)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	groupResp := make([]*dto.Group, 0, len(groups))
	for _, g := range groups {
		groupResp = append(groupResp, &dto.Group{
			ID:    g.ID,
			Name:  g.Name,
			Roles: g.Roles,
		})
	}

	nextPage := &commondto.Paging{
		Limit:  req.Paging.Limit,
		Offset: req.Paging.Offset + req.Paging.Limit,
	}

	if int64(len(groups)) <= req.Paging.Limit {
		nextPage = nil
	}

	if l := int64(len(groupResp)); l > req.Paging.Limit {
		groupResp = groupResp[:req.Paging.Limit]
	}

	transportutil.SendJSONResponse(c, &dto.FindGroupResponse{
		Groups:   groupResp,
		NextPage: nextPage,
	})
}
