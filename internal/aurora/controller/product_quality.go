package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/dto"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/product_quality"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/interceptor"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

type ProductQualityController interface {
	CreateProductQuality(c *gin.Context)
	EditProductQuality(c *gin.Context)
	DeleteProductQuality(c *gin.Context)
	FindProductQuality(c *gin.Context)
}

type productQualityController struct {
	productQualityService product_quality.Service
}

func (s productQualityController) CreateProductQuality(c *gin.Context) {
	req := &dto.CreateProductQualityRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	userID := interceptor.UserIDFromCtx(c)

	id, err := s.productQualityService.CreateProductQuality(c, &product_quality.CreateProductQualityOpts{
		ProductionOrderID: req.ProductionOrderID,
		ProductID:         req.ProductID,
		DeviceID:          req.DeviceID,
		DefectType:        req.DefectType,
		DefectCode:        req.DefectCode,
		DefectLevel:       req.DefectLevel,
		ProductionStageID: req.ProductionStageID,
		DefectiveQuantity: req.DefectiveQuantity,
		GoodQuantity:      req.GoodQuantity,
		Description:       req.Description,
		CreatedBy:         userID,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.CreateProductQualityResponse{
		ID: id,
	})
}

func (s productQualityController) EditProductQuality(c *gin.Context) {
	req := &dto.EditProductQualityRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.productQualityService.EditProductQuality(c, &product_quality.EditProductQualityOpts{
		ID:                req.ID,
		DefectType:        req.DefectType,
		DeviceID:          req.DeviceID,
		DefectCode:        req.DefectCode,
		DefectLevel:       req.DefectLevel,
		ProductionStageID: req.ProductionStageID,
		DefectiveQuantity: req.DefectiveQuantity,
		GoodQuantity:      req.GoodQuantity,
		Description:       req.Description,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.EditProductQualityResponse{})
}

func (s productQualityController) DeleteProductQuality(c *gin.Context) {
	req := &dto.DeleteProductQualityRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.productQualityService.Delete(c, req.ID)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.DeleteProductQualityResponse{})
}

func (s productQualityController) FindProductQuality(c *gin.Context) {
	req := &dto.FindProductQualityRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	productQualitys, cnt, analysis, err := s.productQualityService.FindProductQuality(c, &product_quality.FindProductQualityOpts{
		ProductionOrderID: req.Filter.ProductionOrderID,
		DeviceID: 		   req.Filter.DeviceID,
		DefectType:        req.Filter.DefectType,
		DefectCode:        req.Filter.DefectCode,
		CreatedAtFrom:     req.Filter.CreatedAtFrom,
		CreatedAtTo:       req.Filter.CreatedAtTo,
	}, &repository.Sort{
		Order: repository.SortOrderDESC,
		By:    "ID",
	}, req.Paging.Limit, req.Paging.Offset)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	productQualityResp := make([]*dto.ProductQuality, 0, len(productQualitys))
	for _, f := range productQualitys {
		productQualityResp = append(productQualityResp, toProductQualityResp(f))
	}

	analysisResp := make([]*dto.ProductQualityAnalysis, 0, len(analysis))
	for _, a := range analysis {
		analysisResp = append(analysisResp, &dto.ProductQualityAnalysis{
			DefectType: a.DefectType,
			Count:      a.Count,
		})
	}

	transportutil.SendJSONResponse(c, &dto.FindProductQualityResponse{
		ProductQuality: productQualityResp,
		Total:          cnt.Count,
		Analysis:       analysisResp,
	})
}

func toProductQualityResp(f *product_quality.Data) *dto.ProductQuality {
	return &dto.ProductQuality{
		ID:                  f.ID,
		ProductionOrderID:   f.ProductionOrderID,
		ProductionOrderName: f.ProductionOrderName,
		DeviceID:            f.DeviceID.String,
		ProductID:           f.ProductID.String,
		DefectType:          f.DefectType.String,
		DefectCode:          f.DefectCode.String,
		DefectLevel:         f.DefectLevel,
		ProductionStageID:   f.ProductionStageID.String,
		DefectiveQuantity:   f.DefectiveQuantity,
		GoodQuantity:        f.GoodQuantity,
		Description:         f.Description.String,
		CreatedBy:           f.CreatedBy,
		CreatedAt:           f.CreatedAt,
		UpdatedAt:           f.UpdatedAt,
	}
}

func RegisterProductQualityController(
	r *gin.RouterGroup,
	productQualityService product_quality.Service,
) {
	g := r.Group("product-quality")

	var c ProductQualityController = &productQualityController{
		productQualityService: productQualityService,
	}

	routeutil.AddEndpoint(
		g,
		"create",
		c.CreateProductQuality,
		&dto.CreateProductQualityRequest{},
		&dto.CreateProductQualityResponse{},
		"Create productQuality",
	)

	routeutil.AddEndpoint(
		g,
		"edit",
		c.EditProductQuality,
		&dto.EditProductQualityRequest{},
		&dto.EditProductQualityResponse{},
		"Edit productQuality",
	)

	routeutil.AddEndpoint(
		g,
		"delete",
		c.DeleteProductQuality,
		&dto.DeleteProductQualityRequest{},
		&dto.DeleteProductQualityResponse{},
		"delete productQuality",
	)

	routeutil.AddEndpoint(
		g,
		"find",
		c.FindProductQuality,
		&dto.FindProductQualityRequest{},
		&dto.FindProductQualityResponse{},
		"Find productQuality",
	)
}
