package controller

import (
	"github.com/gin-gonic/gin"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/dto"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

type StatisticsController interface {
	GetStatistics(c *gin.Context)
}

type statisticsController struct {
	statisticsService service.StatisticsService
}

func (s statisticsController) GetStatistics(c *gin.Context) {
	req := &dto.StatisticsRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	res, err := s.statisticsService.GetStatistics(c, req)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, res)
}

func RegisterStatisticsController(
	r *gin.RouterGroup,
	statisticsService service.StatisticsService,
) {
	g := r.Group("statistics")

	var c StatisticsController = &statisticsController{
		statisticsService: statisticsService,
	}

	routeutil.AddEndpoint(
		g,
		"get",
		c.GetStatistics,
		&dto.StatisticsRequest{},
		&dto.StatisticsResponse{},
		"Get statistics",
		routeutil.RegisterOptionSkipAuth,
	)
}
