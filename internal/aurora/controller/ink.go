package controller

import (
	"github.com/gin-gonic/gin"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/ink_export"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/ink_import"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/ink_return"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/interceptor"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/dto"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/ink"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

type InkController interface {
	FindInk(c *gin.Context)
	EditInk(c *gin.Context)
	ImportInk(c *gin.Context)
	FindInkImport(c *gin.Context)
	ExportInk(c *gin.Context)
	FindInkExportByPO(c *gin.Context)
	FindInkExport(c *gin.Context)
	ReturnInk(c *gin.Context)
	FindInkReturn(c *gin.Context)
}

type inkController struct {
	inkService       ink.Service
	inkImportService ink_import.Service
	inkExportService ink_export.Service
	inkReturnService ink_return.Service
}

func (s inkController) ReturnInk(c *gin.Context) {
	req := &dto.CreateInkReturnRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	userID := interceptor.UserIDFromCtx(c)
	inkReturnDetail := make([]*ink_return.CreateInkReturnDetailOpts, 0, len(req.InkReturnDetail))
	for _, f := range req.InkReturnDetail {
		inkReturnDetail = append(inkReturnDetail, &ink_return.CreateInkReturnDetailOpts{
			InkID:       f.InkID,
			InkExportID: f.InkExportID,
			Quantity:    f.Quantity,
			Description: f.Description,
			Data:        f.Data,
		})
	}
	id, err := s.inkReturnService.Create(c, &ink_return.CreateInkReturnOpts{
		Name:            req.Name,
		Code:            req.Code,
		Description:     req.Description,
		Data:            req.Data,
		InkReturnDetail: inkReturnDetail,
		CreatedBy:       userID,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}
	transportutil.SendJSONResponse(c, &dto.CreateInkExportResponse{
		ID: id,
	})
}

func (s inkController) FindInkReturn(c *gin.Context) {
	req := &dto.FindInkReturnsRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}
	inkReturns, cnt, err := s.inkReturnService.Find(c, &ink_return.FindInkReturnOpts{
		Name: req.Filter.Name,
	}, &repository.Sort{
		Order: repository.SortOrderDESC,
		By:    "ID",
	}, req.Paging.Limit, req.Paging.Offset)

	if err != nil {
		transportutil.Error(c, err)
		return
	}
	inkReturnResp := make([]*dto.InkReturn, 0, len(inkReturns))
	for _, f := range inkReturns {
		inkReturnResp = append(inkReturnResp, toInkReturnResp(f))
	}
	transportutil.SendJSONResponse(c, &dto.FindInkReturnsResponse{
		InkReturn: inkReturnResp,
		Total:     cnt.Count,
	})
}

func toInkReturnResp(f *ink_return.InkReturnData) *dto.InkReturn {
	inkReturnDetail := make([]*dto.InkReturnDetail, 0, len(f.InkReturnDetail))
	for _, k := range f.InkReturnDetail {
		d := k.InkData
		inkData := &dto.InkDataExportDetail{
			Name:         d.Name,
			Code:         d.Code,
			ProductCodes: d.ProductCodes,
			Position:     d.Position,
			Location:     d.Location,
			Manufacturer: d.Manufacturer,
			ColorDetail:  d.ColorDetail,
		}

		inkReturnDetail = append(inkReturnDetail, &dto.InkReturnDetail{
			ID:          k.ID,
			InkID:       k.InkID,
			InkData:     inkData,
			InkExportID: k.InkExportID,
			Quantity:    k.Quantity,
			Description: k.Description,
			Data:        k.Data,
		})
	}
	return &dto.InkReturn{
		ID:              f.ID,
		Name:            f.Name,
		Code:            f.Code,
		Description:     f.Description,
		Data:            f.Data,
		InkReturnDetail: inkReturnDetail,
	}
}

func (s inkController) FindInkExport(c *gin.Context) {
	req := &dto.FindInkExportsRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}
	inkExports, cnt, err := s.inkExportService.Find(c, &ink_export.FindInkExportOpts{
		Code:        req.Filter.InkCode,
		ProductName: req.Filter.ProductName,
	}, &repository.Sort{
		Order: repository.SortOrderDESC,
		By:    "ID",
	}, req.Paging.Limit, req.Paging.Offset)

	if err != nil {
		transportutil.Error(c, err)
		return
	}
	inkExportResp := make([]*dto.InkExport, 0, len(inkExports))
	for _, f := range inkExports {
		inkExportResp = append(inkExportResp, toInkExportResp(f))
	}
	transportutil.SendJSONResponse(c, &dto.FindInkExportsResponse{
		InkExport: inkExportResp,
		Total:     cnt.Count,
	})
}

func toInkExportResp(f *ink_export.InkExportData) *dto.InkExport {
	inkExportDetail := make([]*dto.InkExportDetail, 0, len(f.InkExportDetail))
	for _, k := range f.InkExportDetail {
		d := k.InkData
		inkData := &dto.InkDataExportDetail{
			Name:         d.Name,
			Code:         d.Code,
			ProductCodes: d.ProductCodes,
			Position:     d.Position,
			Location:     d.Location,
			Manufacturer: d.Manufacturer,
			ColorDetail:  d.ColorDetail,
		}

		inkExportDetail = append(inkExportDetail, &dto.InkExportDetail{
			ID:          k.ID,
			InkID:       k.InkID,
			InkData:     inkData,
			Quantity:    k.Quantity,
			Description: k.Description,
			Data:        k.Data,
		})
	}

	po := f.ProductionOrderData
	productionOrderData := &dto.ProductionOrderData{}
	if po != nil {
		productionOrderData = &dto.ProductionOrderData{
			ID:          po.ID,
			Name:        po.Name,
			ProductCode: po.ProductCode,
			ProductName: po.ProductName,
		}
	}

	return &dto.InkExport{
		ID:                  f.ID,
		Name:                f.Name,
		Code:                f.Code,
		ProductionOrderID:   f.ProductionOrderID,
		Description:         f.Description,
		Data:                f.Data,
		Status:              f.Status,
		CreatedBy:           f.CreatedBy,
		CreatedAt:           f.CreatedAt,
		InkExportDetail:     inkExportDetail,
		ProductionOrderData: productionOrderData,
	}

}

func (s inkController) ExportInk(c *gin.Context) {
	req := &dto.CreateInkExportRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}
	userID := interceptor.UserIDFromCtx(c)
	inkExportDetail := make([]*ink_export.CreateInkExportDetailOpts, 0, len(req.InkExportDetail))
	for _, f := range req.InkExportDetail {
		inkExportDetail = append(inkExportDetail, &ink_export.CreateInkExportDetailOpts{
			InkID:       f.InkID,
			Quantity:    f.Quantity,
			Description: f.Description,
			Data:        f.Data,
		})
	}
	id, err := s.inkExportService.Create(c, &ink_export.CreateInkExportOpts{
		Name:              req.Name,
		Code:              req.Code,
		ProductionOrderID: req.ProductionOrderID,
		ExportDate:        req.ExportDate,
		Description:       req.Description,
		Data:              req.Data,
		CreatedBy:         userID,
		InkExportDetail:   inkExportDetail,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}
	transportutil.SendJSONResponse(c, &dto.CreateInkExportResponse{
		ID: id,
	})
}

func (s inkController) FindInkExportByPO(c *gin.Context) {
	req := &dto.FindInkExportByPORequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}
	exportDetails, err := s.inkExportService.FindImportDetailByPOID(c, req.ProductionOrderID)
	if err != nil {
		transportutil.Error(c, err)
		return
	}
	exportDetailResp := make([]*dto.InkExportDetail, 0, len(exportDetails))
	for _, f := range exportDetails {
		d := f.InkData
		inkData := &dto.InkDataExportDetail{
			Name:         d.Name,
			Code:         d.Code,
			ProductCodes: d.ProductCodes,
			Position:     d.Position,
			Location:     d.Location,
			Manufacturer: d.Manufacturer,
			ColorDetail:  d.ColorDetail,
		}
		exportDetailResp = append(exportDetailResp, &dto.InkExportDetail{
			ID:          f.ID,
			InkID:       f.InkID,
			InkExportID: f.InkExportID,
			InkData:     inkData,
			Quantity:    f.Quantity,
			Description: f.Description,
			Data:        f.Data,
		})
	}
	transportutil.SendJSONResponse(c, &dto.FindInkExportByPOResponse{
		InkExportDetail: exportDetailResp,
	})
}

func (s inkController) FindInkImport(c *gin.Context) {
	req := &dto.FindInkImportsRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	imports, cnt, err := s.inkImportService.Find(c, &ink_import.FindInkImportOpts{
		Name:   req.Filter.Name,
		ID:     req.Filter.ID,
		Status: req.Filter.Status,
	}, &repository.Sort{
		Order: repository.SortOrderDESC,
		By:    "ID",
	}, req.Paging.Limit, req.Paging.Offset)

	if err != nil {
		transportutil.Error(c, err)
		return
	}

	importResp := make([]*dto.InkImport, 0, len(imports))
	for _, f := range imports {
		importResp = append(importResp, toInkImportResp(f))
	}

	transportutil.SendJSONResponse(c, &dto.FindInkImportsResponse{
		InkImport: importResp,
		Total:     cnt.Count,
	})
}

func toInkImportDetailResp(f []*ink_import.InkImportDetail) []*dto.InkImportDetail {
	importDetailResp := make([]*dto.InkImportDetail, 0, len(f))
	for _, f := range f {
		importDetailResp = append(importDetailResp, &dto.InkImportDetail{
			ID:             f.ID,
			InkID:          f.InkID,
			Name:           f.Name,
			Code:           f.Code,
			ProductCodes:   f.ProductCodes,
			Position:       f.Position,
			Location:       f.Location,
			Manufacturer:   f.Manufacturer,
			ColorDetail:    f.ColorDetail,
			Quantity:       f.Quantity,
			ExpirationDate: f.ExpirationDate,
			Description:    f.Description,
			Data:           f.Data,
		})
	}
	return importDetailResp
}

func toInkImportResp(f *ink_import.InkImportData) *dto.InkImport {
	return &dto.InkImport{
		ID:              f.ID,
		Name:            f.Name,
		Code:            f.Code,
		Description:     f.Description,
		Data:            f.Data,
		Status:          f.Status,
		InkImportDetail: toInkImportDetailResp(f.InkImportDetail),
		CreatedBy:       f.CreatedBy,
		CreatedAt:       f.CreatedAt,
	}
}

func (s inkController) ImportInk(c *gin.Context) {
	req := &dto.CreateInkImportRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}
	userID := interceptor.UserIDFromCtx(c)
	inkImportDetail := make([]*ink_import.CreateInkImportDetailOpts, 0, len(req.InkImportDetail))
	for _, f := range req.InkImportDetail {
		inkImportDetail = append(inkImportDetail, &ink_import.CreateInkImportDetailOpts{
			Name:           f.Name,
			Code:           f.Code,
			ProductCodes:   f.ProductCodes,
			Position:       f.Position,
			Location:       f.Location,
			Manufacturer:   f.Manufacturer,
			ColorDetail:    f.ColorDetail,
			Quantity:       f.Quantity,
			ExpirationDate: f.ExpirationDate,
			Description:    f.Description,
			Data:           f.Data,
		})
	}
	id, err := s.inkImportService.Create(c, &ink_import.CreateInkImportOpts{
		Name:            req.Name,
		Code:            req.Code,
		Description:     req.Description,
		Data:            req.Data,
		InkImportDetail: inkImportDetail,
		CreatedBy:       userID,
	})

	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.CreateInkImportResponse{
		ID: id,
	})
}

func (s inkController) FindInk(c *gin.Context) {
	req := &dto.FindInkRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	inks, cnt, err := s.inkService.Find(c, &ink.FindInkOpts{
		Name:   req.Filter.Name,
		ID:     req.Filter.ID,
		Code:   req.Filter.Code,
		Status: req.Filter.Status,
	}, &repository.Sort{
		Order: repository.SortOrderDESC,
		By:    "ID",
	}, req.Paging.Limit, req.Paging.Offset)

	if err != nil {
		transportutil.Error(c, err)
		return
	}

	inkResp := make([]*dto.Ink, 0, len(inks))
	for _, f := range inks {
		inkResp = append(inkResp, toInkResp(f))
	}

	transportutil.SendJSONResponse(c, &dto.FindInksResponse{
		Ink:   inkResp,
		Total: cnt.Count,
	})
}
func (s inkController) EditInk(c *gin.Context) {
	req := &dto.EditInkRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}
	userId := interceptor.UserIDFromCtx(c)
	err = s.inkService.Edit(c, &ink.EditInkOpts{
		ID:             req.ID,
		Name:           req.Name,
		Code:           req.Code,
		ProductCodes:   req.ProductCodes,
		Position:       req.Position,
		Location:       req.Location,
		Manufacturer:   req.Manufacturer,
		ColorDetail:    req.ColorDetail,
		ExpirationDate: req.ExpirationDate,
		Description:    req.Description,
		Data:           req.Data,
		Status:         req.Status,
		Quantity:       req.Quantity,
		UpdatedBy:      userId,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}
	transportutil.SendJSONResponse(c, &dto.EditInkResponse{})
}

func toInkResp(f *ink.InkData) *dto.Ink {
	return &dto.Ink{
		ID:             f.ID,
		ImportID:       f.ImportID.String,
		Name:           f.Name,
		Code:           f.Code,
		ProductCodes:   f.ProductCodes,
		Position:       f.Position,
		Location:       f.Location,
		Manufacturer:   f.Manufacturer,
		ColorDetail:    f.ColorDetail,
		Quantity:       f.Quantity,
		ExpirationDate: f.ExpirationDate,
		Description:    f.Description.String,
		Data:           f.Data,
		Status:         f.Status,
		CreatedBy:      f.CreatedBy,
		UpdatedBy:      f.UpdatedBy,
		CreatedAt:      f.CreatedAt,
		UpdatedAt:      f.UpdatedAt,
	}
}

func RegisterInkController(
	r *gin.RouterGroup,
	inkService ink.Service,
	inkImportService ink_import.Service,
	inkReturnService ink_return.Service,
	inkExportService ink_export.Service,
) {
	g := r.Group("ink")

	var c InkController = &inkController{
		inkService:       inkService,
		inkReturnService: inkReturnService,
		inkImportService: inkImportService,
		inkExportService: inkExportService,
	}

	routeutil.AddEndpoint(
		g,
		"find",
		c.FindInk,
		&dto.FindInkRequest{},
		&dto.FindInksResponse{},
		"Find ink",
	)

	routeutil.AddEndpoint(
		g,
		"edit",
		c.EditInk,
		&dto.EditInkRequest{},
		&dto.EditInkRequest{},
		"Edit ink",
	)

	routeutil.AddEndpoint(
		g,
		"import",
		c.ImportInk,
		&dto.CreateInkImportRequest{},
		&dto.CreateInkImportResponse{},
		"Import ink",
	)

	routeutil.AddEndpoint(
		g,
		"find-ink-import",
		c.FindInkImport,
		&dto.FindInkImportsRequest{},
		&dto.FindInkImportsResponse{},
		"Find import ink",
	)

	routeutil.AddEndpoint(
		g,
		"export",
		c.ExportInk,
		&dto.CreateInkExportRequest{},
		&dto.CreateInkExportResponse{},
		"export ink",
	)

	routeutil.AddEndpoint(
		g,
		"find-ink-export",
		c.FindInkExport,
		&dto.FindInkExportsRequest{},
		&dto.FindInkExportsResponse{},
		"Find export ink",
	)
	routeutil.AddEndpoint(
		g,
		"find-ink-export-detail-by-po",
		c.FindInkExportByPO,
		&dto.FindInkExportByPORequest{},
		&dto.FindInkExportByPOResponse{},
		"Find export detail ink by production order",
	)

	routeutil.AddEndpoint(
		g,
		"return",
		c.ReturnInk,
		&dto.CreateInkReturnRequest{},
		&dto.CreateInkReturnResponse{},
		"return ink",
	)

	routeutil.AddEndpoint(
		g,
		"find-ink-return",
		c.FindInkReturn,
		&dto.FindInkReturnsRequest{},
		&dto.FindInkReturnsResponse{},
		"Find return ink",
	)
}
