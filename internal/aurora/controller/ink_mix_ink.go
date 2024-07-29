package controller

import (
	"github.com/gin-gonic/gin"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/dto"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/ink"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/interceptor"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

func (s inkController) MixInk(c *gin.Context) {
	req := &dto.CreateInkMixingRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}
	userID := interceptor.UserIDFromCtx(c)
	inkFormulation := make([]ink.CreateFormulation, 0)
	for _, inkForm := range req.InkFormulation {
		inkFormulation = append(inkFormulation, ink.CreateFormulation{
			InkID:       inkForm.InkID,
			Quantity:    inkForm.Quantity,
			Description: inkForm.Description,
		})
	}
	// Create ink mixing
	inkMixingData := &ink.CreateInkMixingOpts{
		Name:           req.Name,
		Code:           req.Code,
		ProductCodes:   req.ProductCodes,
		Quantity:       req.Quantity,
		ExpirationDate: req.ExpirationDate,
		Position:       req.Position,
		Location:       req.Location,
		Description:    req.Description,
		InkFormula:     inkFormulation,
		Status:         enum.CommonStatusActive,
		CreatedBy:      userID,
	}
	inkMixingID, err := s.inkService.MixInk(c, inkMixingData)
	if err != nil {
		transportutil.Error(c, err)
		return
	}
	transportutil.SendJSONResponse(c, &dto.CreateInkMixingResponse{
		ID: inkMixingID,
	})
}

func (s inkController) EditInkMixing(c *gin.Context) {
	req := &dto.EditInkMixingRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}
	userID := interceptor.UserIDFromCtx(c)
	inkFormulation := make([]ink.EditInkFormulation, 0)
	for _, inkForm := range req.InkFormulation {
		inkFormulation = append(inkFormulation, ink.EditInkFormulation{
			ID:          inkForm.ID,
			InkID:       inkForm.InkID,
			Quantity:    inkForm.Quantity,
			Description: inkForm.Description,
		})
	}
	// Create ink mixing
	inkMixingData := &ink.EditInkMixingOpts{
		ID:             req.ID,
		Name:           req.Name,
		Code:           req.Code,
		ProductCodes:   req.ProductCodes,
		Quantity:       req.Quantity,
		ExpirationDate: req.ExpirationDate,
		Position:       req.Position,
		Location:       req.Location,
		Description:    req.Description,
		InkFormula:     inkFormulation,
		Status:         enum.CommonStatusActive,
		UpdatedBy:      userID,
	}
	err := s.inkService.EditInkMixing(c, inkMixingData)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

}

func (s inkController) FindInkMixing(c *gin.Context) {
	req := &dto.FindInkMixingRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}
	inkMixingData, count, err := s.inkService.FindInkMixing(c, &ink.FindInkMixingOpts{
		IDs:    req.Filter.IDs,
		Search: req.Filter.Search,
		InkID:  req.Filter.InkID,
		Limit:  req.Paging.Limit,
		Offset: req.Paging.Offset,
	})

	if err != nil {
		transportutil.Error(c, err)
		return
	}
	results := make([]*dto.MixInk, 0)
	for _, ink := range inkMixingData {
		inkFormula := make([]dto.InkMixingFormulation, 0)
		for _, detail := range ink.InkFormula {
			inkFormula = append(inkFormula, dto.InkMixingFormulation{
				ID:          detail.ID,
				InkID:       detail.InkID,
				Quantity:    detail.Quantity,
				Description: detail.Description,
			})
		}
		results = append(results, &dto.MixInk{
			ID:             ink.ID,
			Name:           ink.Name,
			Code:           ink.Code,
			ProductCodes:   ink.ProductCodes,
			Quantity:       ink.Quantity,
			ExpirationDate: ink.ExpirationDate,
			Position:       ink.Position,
			Location:       ink.Location,
			Description:    ink.Description,
			InkFormulation: inkFormula,
			Status:         ink.Status,
			CreatedBy:      ink.CreatedBy,
			CreatedAt:      ink.CreatedAt,
			UpdatedAt:      ink.UpdatedAt,
			CreatedByName:  ink.CreatedByName,
			UpdatedByName:  ink.UpdatedByName,
		})
	}
	transportutil.SendJSONResponse(c, &dto.FindInkMixingResponse{
		MixInk: results,
		Total:  count.Count,
	})
}
