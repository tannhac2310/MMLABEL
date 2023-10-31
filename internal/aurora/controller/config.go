package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/dto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/constants"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

type ConfigController interface {
	GetAppConfig(c *gin.Context)
}

type configController struct {
}

func (c *configController) GetAppConfig(ctx *gin.Context) {
	transportutil.SendJSONResponse(ctx, &dto.AppConfigResponse{
		HydraPermissions: constants.HydraPermissions,
		GezuPermissions:  constants.GezuPermissionsList,
	})
}

func RegisterConfigController(r *gin.RouterGroup) {
	g := r.Group("app")

	c := &configController{}
	routeutil.AddEndpoint(
		g,
		"config",
		c.GetAppConfig,
		&dto.AppConfigRequest{},
		&dto.AppConfigResponse{},
		"get app config",
	)
}
