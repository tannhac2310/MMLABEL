package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	pkgConfig "mmlabel.gitlab.com/mm-printing-backend/pkg/configs"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/role"

	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/dto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/casbinrule"
)

type RuleController interface {
	CreateRule(c *gin.Context)
	FindRules(c *gin.Context)
	Delete(c *gin.Context)
}

type ruleController struct {
	ruleService casbinrule.Service
	redisDB     redis.Cmdable
	baseCfg     *pkgConfig.BaseConfig
	roleService role.Service
}

func RegisterRuleController(
	r *gin.RouterGroup,
	ruleService casbinrule.Service,
	redisDB redis.Cmdable,
	baseCfg *pkgConfig.BaseConfig,
	roleService role.Service,

) {
	g := r.Group("rule")

	var c RuleController = &ruleController{
		ruleService: ruleService,
		redisDB:     redisDB,
		baseCfg:     baseCfg,
		roleService: roleService,
	}

	routeutil.AddEndpoint(
		g,
		"create-rule",
		c.CreateRule,
		&dto.CreateRuleRequest{},
		&dto.CreateRuleResponse{},
		"Create rule",
	)

	routeutil.AddEndpoint(
		g,
		"delete-rule",
		c.Delete,
		&dto.DeleteRuleRequest{},
		&dto.DeleteRuleResponse{},
		"Delete one rule",
	)
}
