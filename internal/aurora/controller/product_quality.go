package controller

import (
	"fmt"

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
	inspectionErrors := make([]*product_quality.InspectionError, 0, len(req.InspectionErrors))
	for _, e := range req.InspectionErrors {
		if e.DeviceID == "" {
			transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage("DeviceID is required"))
			return
		}
		inspectionErrors = append(inspectionErrors, &product_quality.InspectionError{
			DeviceID:         e.DeviceID,
			DeviceName:       e.DeviceName,
			ErrorType:        e.ErrorType,
			Quantity:         e.Quantity,
			NhanVienThucHien: e.NhanVienThucHien,
			Note:             e.Note,
		})
	}
	id, err := s.productQualityService.CreateProductQuality(c, &product_quality.CreateProductQualityOpts{
		ProductionOrderID:   req.ProductionOrderID,
		InspectionDate:      req.InspectionDate,
		InspectorName:       req.InspectorName,
		Quantity:            req.Quantity,
		Note:                req.Note,
		ProductID:           req.ProductID,
		SoLuongHopDong:      req.SoLuongHopDong,
		SoLuongIn:           req.SoLuongIn,
		NguoiKiemTra:        req.NguoiKiemTra,
		NguoiPheDuyet:       req.NguoiPheDuyet,
		SoLuongThanhPhamDat: req.SoLuongThanhPhamDat,
		CreatedBy:           userID,
		InspectionErrors:    inspectionErrors,
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
	inspectionErrors := make([]*product_quality.EditInspectionError, 0, len(req.InspectionErrors))
	for _, e := range req.InspectionErrors {
		inspectionErrors = append(inspectionErrors, &product_quality.EditInspectionError{
			DeviceID:         e.DeviceID,
			DeviceName:       e.DeviceName,
			ErrorType:        e.ErrorType,
			Quantity:         e.Quantity,
			Note:             e.Note,
			NhanVienThucHien: e.NhanVienThucHien,
		})
	}
	err = s.productQualityService.EditProductQuality(c, &product_quality.EditProductQualityOpts{
		ID:                  req.ID,
		ProductionOrderID:   req.ProductionOrderID,
		InspectionDate:      req.InspectionDate,
		InspectorName:       req.InspectorName,
		Quantity:            req.Quantity,
		Note:                req.Note,
		ProductID:           req.ProductID,
		SoLuongHopDong:      req.SoLuongHopDong,
		SoLuongIn:           req.SoLuongIn,
		NguoiKiemTra:        req.NguoiKiemTra,
		NguoiPheDuyet:       req.NguoiPheDuyet,
		SoLuongThanhPhamDat: req.SoLuongThanhPhamDat,
		InspectionErrors:    inspectionErrors,
		CreatedBy:           interceptor.UserIDFromCtx(c),
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

	productQuality, cnt, err := s.productQualityService.FindProductQuality(c, &product_quality.FindProductQualityOpts{
		IDs:               req.Filter.IDs,
		ProductionOrderID: req.Filter.ProductionOrderID,
		DefectTypes:       req.Filter.DefectTypes,
		CreatedAtFrom:     req.Filter.CreatedAtFrom,
		CreatedAtTo:       req.Filter.CreatedAtTo,
	}, &repository.Sort{
		Order: repository.SortOrderDESC,
		By:    "created_at",
	}, req.Paging.Limit, req.Paging.Offset)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	productQualityResp := make([]*dto.ProductQuality, 0, len(productQuality))
	for _, f := range productQuality {
		productQualityResp = append(productQualityResp, toProductQualityResp(f))
	}
	fmt.Println(cnt)
	transportutil.SendJSONResponse(c, &dto.FindProductQualityResponse{
		ProductQuality: productQualityResp,
		Total:          cnt.Count,
		Analysis:       nil,
	})
}
func toProductQualityResp(f *product_quality.Data) *dto.ProductQuality {
	inspectionErrors := make([]dto.InspectionError, 0)
	for _, err := range f.InspectionErrors {
		inspectionErrors = append(inspectionErrors, dto.InspectionError{
			ID:               err.ID,
			DeviceID:         err.DeviceID,
			DeviceName:       err.DeviceName,
			InspectionFormID: err.InspectionFormID,
			ErrorType:        err.ErrorType,
			Quantity:         err.Quantity,
			Note:             err.Note,
			NhanVienThucHien: err.NhanVienThucHien,
		})
	}

	return &dto.ProductQuality{
		ID:                  f.ID,
		ProductionOrderID:   f.ProductionOrderID,
		ProductionOrderCode: f.ProductionOrderCode,
		ProductionOrderName: f.ProductionOrderName,
		InspectionDate:      f.InspectionDate,
		InspectorName:       f.InspectorName,
		Quantity:            f.Quantity,
		ProductID:           f.ProductID,
		ProductCode:         f.ProductCode,
		ProductName:         f.ProductName,
		CustomerID:          f.CustomerID,
		CustomerName:        f.CustomerName,
		CustomerCode:        f.CustomerCode,
		SoLuongHopDong:      f.SoLuongHopDong,
		SoLuongIn:           f.SoLuongIn,
		MaDonDatHang:        f.MaDonDatHang,
		OrderData: struct {
			ID          string `json:"id"`
			MaDatHangMm string `json:"maDatHangMm"`
			Status      string `json:"status"`
		}{ID: f.OrderID, MaDatHangMm: f.MaDonDatHang, Status: f.TrangThaiDonHang},
		NguoiKiemTra:        f.NguoiKiemTra,
		NguoiPheDuyet:       f.NguoiPheDuyet,
		SoLuongThanhPhamDat: f.SoLuongThanhPhamDat,
		Note:                f.Note,
		InspectionErrors:    inspectionErrors,
		CreatedBy:           f.CreatedBy,
		UpdatedBy:           f.UpdatedBy,
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
