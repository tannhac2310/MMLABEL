package controller

import (
	"github.com/gin-gonic/gin"
	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/dto"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

func (b *ruleController) FindRules(c *gin.Context) {
	req := &dto.FindRulesRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	rules, total, err := b.ruleService.FindCasbinRules(
		c,
		req.Paging.Limit+1,
		req.Paging.Offset)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	ruleResp := make([]*dto.Rule, 0, len(rules))
	for _, f := range rules {
		ruleResp = append(ruleResp, &dto.Rule{
			Role: f.V0,
			Rule: f.V1,
		})
	}

	nextPage := &commondto.Paging{
		Limit:  req.Paging.Limit,
		Offset: req.Paging.Offset + req.Paging.Limit,
	}

	if int64(len(rules)) <= req.Paging.Limit {
		nextPage = nil
	}

	if l := int64(len(ruleResp)); l > req.Paging.Limit {
		ruleResp = ruleResp[:req.Paging.Limit]
	}

	transportutil.SendJSONResponse(c, &dto.FindRulesResponse{
		Rules:    ruleResp,
		NextPage: nextPage,
		Total:    total.Count,
	})
}
