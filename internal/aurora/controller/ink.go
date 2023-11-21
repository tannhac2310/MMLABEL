package controller

import (
	"github.com/gin-gonic/gin"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/ink_export"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/ink_import"
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
	ImportInk(c *gin.Context)
	FindInkImport(c *gin.Context)
	ExportInk(c *gin.Context)
	FindInkExport(c *gin.Context)
}

type inkController struct {
	inkService       ink.Service
	inkImportService ink_import.Service
	inkExportService ink_export.Service
}

func (s inkController) FindInkExport(c *gin.Context) {
	req := &dto.FindInkExportsRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}
	inkExports, cnt, err := s.inkExportService.Find(c, &ink_export.FindInkExportOpts{
		Name: req.Filter.Name,
	}, &repository.Sort{
		Order: repository.SortOrderASC,
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
	return &dto.InkExport{
		ID:                f.ID,
		Name:              f.Name,
		Code:              f.Code,
		ProductionOrderID: f.ProductionOrderID,
		Description:       f.Description,
		Data:              f.Data,
		Status:            f.Status,
		CreatedBy:         f.CreatedBy,
		CreatedAt:         f.CreatedAt,
		InkExportDetail:   inkExportDetail,
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
		Order: repository.SortOrderASC,
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
		Status: req.Filter.Status,
	}, &repository.Sort{
		Order: repository.SortOrderASC,
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
	inkExportService ink_export.Service,
) {
	g := r.Group("ink")

	var c InkController = &inkController{
		inkService:       inkService,
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
}
