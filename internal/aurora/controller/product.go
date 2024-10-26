package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/dto"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/product"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/interceptor"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

type ProductController interface {
	CreateProduct(ctx *gin.Context)
	UpdateProduct(ctx *gin.Context)
	DeleteProduct(ctx *gin.Context)
	FindProduct(ctx *gin.Context)
}

type productController struct {
	productService product.Service
}

func (p productController) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	userID := interceptor.UserIDFromCtx(c)
	uf := make([]*product.UserField, 0)
	for _, f := range req.UserField {
		uf = append(uf, &product.UserField{
			Key:   f.Key,
			Value: f.Value,
		})
	}
	id, err := p.productService.CreateProduct(c, &product.CreateProductOpts{
		Name:        req.Name,
		Code:        req.Code,
		CustomerID:  req.CustomerID,
		SaleID:      req.SaleID,
		Description: req.Description,
		Data:        req.Data,
		CreatedBy:   userID,
		UserField:   uf,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, dto.CreateProductResponse{
		ID: id,
	})
}

func (p productController) UpdateProduct(c *gin.Context) {
	var req dto.UpdateProductOpts
	if err := c.ShouldBindJSON(&req); err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	userID := interceptor.UserIDFromCtx(c)
	uf := make([]*product.UserField, 0)
	for _, f := range req.UserField {
		uf = append(uf, &product.UserField{
			Key:   f.Key,
			Value: f.Value,
		})
	}
	err := p.productService.UpdateProduct(c, &product.UpdateProductOpts{
		ID:          req.ID,
		Name:        req.Name,
		Code:        req.Code,
		CustomerID:  req.CustomerID,
		SaleID:      req.SaleID,
		Description: req.Description,
		Data:        req.Data,
		UpdatedBy:   userID,
		UserField:   uf,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, dto.UpdateProductResponse{})
}

func (p productController) DeleteProduct(c *gin.Context) {
	var req dto.DeleteProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err := p.productService.DeleteProduct(c, req.ID)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, dto.DeleteProductResponse{})
}

func (p productController) FindProduct(c *gin.Context) {
	var req dto.SearchProductFilter
	if err := c.ShouldBindJSON(&req); err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	filter := product.FindProductOpts{
		Name:           req.Filter.Name,
		Code:           req.Filter.Code,
		CustomerID:     req.Filter.CustomerID,
		SaleID:         req.Filter.SaleID,
		ProductPlanID:  req.Filter.ProductionPlanID,
		ProductOrderID: req.Filter.ProductionOrderID,
	}

	sort := &repository.Sort{
		By:    "id",
		Order: "DESC",
	}

	data, count, err := p.productService.FindProduct(c, &filter, sort, req.Paging.Limit, req.Paging.Offset)
	if err != nil {
		transportutil.Error(c, err)
		return
	}
	result := make([]*dto.Product, 0)
	for _, d := range data {

		v := &dto.Product{
			ID:         d.ID,
			Name:       d.Name,
			Code:       d.Code,
			CustomerID: d.CustomerID,

			ProductionPlanID: d.ProductionPlanID,
			SaleID:           d.SaleID,
			Description:      d.Description,
			Data:             d.Data,
			CreatedAt:        d.CreatedAt,
			UpdatedAt:        d.UpdatedAt,
			CreatedBy:        d.CreatedBy,
			CreatedByName:    d.CreatedByName,
			UpdatedBy:        d.UpdatedBy,
			UpdatedByName:    d.UpdatedByName,
		}
		if d.CustomerData != nil {
			v.CustomerData = &dto.Customer{
				ID:   d.CustomerData.ID,
				Name: d.CustomerData.Name,
			}
		}
		if d.ProductionPlanData != nil {
			v.ProductionPlanData = &dto.ProductionPlan{
				ID:   d.ProductionPlanData.ID,
				Name: d.ProductionPlanData.Name,
			}
		}
		result = append(result, v)
	}
	transportutil.SendJSONResponse(c, dto.SearchProductResponse{
		Product: result,
		Total:   count.Count,
	})

}

func RegisterProductController(
	r *gin.RouterGroup, productService product.Service,
) {
	g := r.Group("product")
	var c ProductController = &productController{
		productService: productService,
	}

	routeutil.AddEndpoint(
		g, "create",
		c.CreateProduct,
		&dto.CreateProductRequest{},
		&dto.CreateProductResponse{},
		"create product",
	)

	routeutil.AddEndpoint(
		g, "update",
		c.UpdateProduct,
		&dto.UpdateProductOpts{},
		&dto.UpdateProductResponse{},
		"update product",
	)

	routeutil.AddEndpoint(
		g, "delete",
		c.DeleteProduct,
		&dto.DeleteProductRequest{},
		&dto.DeleteProductResponse{},
		"delete product",
	)

	routeutil.AddEndpoint(
		g, "find",
		c.FindProduct,
		&dto.SearchProductFilter{},
		&dto.SearchProductResponse{},
		"find product",
	)

}
