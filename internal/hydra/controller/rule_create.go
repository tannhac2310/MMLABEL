package controller

import (
	"github.com/gin-gonic/gin"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/interceptor"

	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/dto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

func (b *ruleController) CreateRule(c *gin.Context) {
	req := &dto.CreateRuleRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	_, err = b.roleService.AddPolicy(req.Role, req.Rule)

	if err != nil {
		transportutil.Error(c, err)
		return
	}

	interceptor.ForceTokenInValid(c, b.baseCfg, b.redisDB)
	transportutil.SendJSONResponse(c, &dto.CreateRuleResponse{})
}
