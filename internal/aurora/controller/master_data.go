package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/dto"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/master_data"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/interceptor"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

type MasterDataController interface {
	InsertMasterData(ctx *gin.Context)
	UpdateMasterData(ctx *gin.Context)
	DeleteMasterData(ctx *gin.Context)
	GetMasterData(ctx *gin.Context)
}

type masterDataController struct {
	masterDataService master_data.Service
}

func (m masterDataController) InsertMasterData(ctx *gin.Context) {
	req := &dto.CreateMasterDataRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		transportutil.Error(ctx, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}
	userID := interceptor.UserIDFromCtx(ctx)
	uf := make([]master_data.CreateMasterDataUserField, 0)
	for _, userField := range req.UserFields {
		uf = append(uf, master_data.CreateMasterDataUserField{
			FieldName:  userField.FieldName,
			FieldValue: userField.FieldValue,
		})
	}
	masterData := &master_data.CreateMasterDataOpts{
		Type:        req.Type,
		Name:        req.Name,
		Description: req.Description,
		UserFields:  uf,
		CreatedBy:   userID,
	}
	masterDataID, err := m.masterDataService.CreateMasterData(ctx, masterData)
	if err != nil {
		transportutil.Error(ctx, err)
		return
	}
	transportutil.SendJSONResponse(ctx, &dto.CreateMasterDataResponse{
		ID: masterDataID,
	})
}

func (m masterDataController) UpdateMasterData(ctx *gin.Context) {
	req := &dto.UpdateMasterDataRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		transportutil.Error(ctx, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}
	userID := interceptor.UserIDFromCtx(ctx)
	uf := make([]master_data.CreateMasterDataUserField, 0)
	for _, userField := range req.UserFields {
		uf = append(uf, master_data.CreateMasterDataUserField{
			FieldName:  userField.FieldName,
			FieldValue: userField.FieldValue,
		})
	}
	masterData := &master_data.UpdateMasterDataOpts{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		UserFields:  uf,
		UpdateBy:    userID,
	}
	err := m.masterDataService.UpdateMasterData(ctx, masterData)
	if err != nil {
		transportutil.Error(ctx, err)
		return
	}
	transportutil.SendJSONResponse(ctx, nil)
}

func (m masterDataController) DeleteMasterData(ctx *gin.Context) {
	req := &dto.DeleteMasterDataRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		transportutil.Error(ctx, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}
	err := m.masterDataService.DeleteMasterData(ctx, &master_data.DeleteMasterDataOpts{
		ID: req.ID,
	})
	if err != nil {
		transportutil.Error(ctx, err)
		return
	}
	transportutil.SendJSONResponse(ctx, nil)
}

func (m masterDataController) GetMasterData(ctx *gin.Context) {
	req := &dto.SearchMasterDataRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		transportutil.Error(ctx, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}
	data, total, err := m.masterDataService.FindMasterData(ctx, &master_data.FindMasterDataOpts{
		ID:     req.Filter.ID,
		Type:   req.Filter.Type,
		Limit:  req.Paging.Limit,
		Offset: req.Paging.Offset,
	})
	if err != nil {
		transportutil.Error(ctx, err)
		return
	}

	res := make([]*dto.MasterData, 0, len(data))
	for _, d := range data {
		uf := make([]*dto.MasterDataUserField, 0, len(d.UserFields))
		for _, f := range d.UserFields {
			uf = append(uf, &dto.MasterDataUserField{
				ID:           f.ID,
				MasterDataID: f.MasterDataID,
				FieldName:    f.FieldName,
				FieldValue:   f.FieldValue,
			})
		}
		res = append(res, &dto.MasterData{
			ID:          d.ID,
			Type:        d.Type,
			Name:        d.Name,
			Description: d.Description,
			UserFields:  uf,
			CreatedBy:   d.CreatedBy,
			CreatedAt:   d.CreatedAt,
		})
	}

	transportutil.SendJSONResponse(ctx, &dto.GetMasterDataResponse{
		MasterData: res,
		Total:      total,
	})
}

func NewMasterDataController(masterDataService master_data.Service) MasterDataController {
	return &masterDataController{
		masterDataService: masterDataService,
	}
}

func RegisterMasterDataController(r *gin.RouterGroup, masterDataService master_data.Service) {
	g := r.Group("master-data")

	var c MasterDataController = &masterDataController{
		masterDataService: masterDataService,
	}

	routeutil.AddEndpoint(
		g,
		"create",
		c.InsertMasterData,
		&dto.CreateMasterDataRequest{},
		&dto.CreateMasterDataResponse{},
		"Create master data",
	)

	routeutil.AddEndpoint(
		g,
		"update",
		c.UpdateMasterData,
		&dto.UpdateMasterDataRequest{},
		&dto.UpdateMasterDataResponse{},
		"Update master data",
	)

	routeutil.AddEndpoint(
		g,
		"delete",
		c.DeleteMasterData,
		&dto.DeleteMasterDataRequest{},
		&dto.DeleteMasterDataResponse{},
		"Delete master data",
	)

	routeutil.AddEndpoint(
		g,
		"find-master-data",
		c.GetMasterData,
		&dto.SearchMasterDataRequest{},
		&dto.GetMasterDataResponse{},
		"Get master data",
	)
}
