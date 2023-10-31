package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/dto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/group"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

func (g *groupController) CreateGroup(c *gin.Context) {

	req := &dto.CreateGroupRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = g.groupService.CreateGroup(c, &group.CreateGroupOpts{
		ID:    req.ID,
		Name:  req.Name,
		Roles: req.Roles,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.CreateGroupResponse{
		ID: req.ID,
	})
}
