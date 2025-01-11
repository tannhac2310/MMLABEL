package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/dto"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/device_config"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/interceptor"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

type DeviceConfigController interface {
	CreateDeviceConfig(c *gin.Context)
	EditDeviceConfig(c *gin.Context)
	DeleteDeviceConfig(c *gin.Context)
	FindDeviceConfigs(c *gin.Context)
}

type deviceConfigController struct {
	deviceConfigService device_config.Service
}

func (s deviceConfigController) CreateDeviceConfig(c *gin.Context) {
	req := &dto.CreateDeviceConfigRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}
	userId := interceptor.UserIDFromCtx(c)
	id, err := s.deviceConfigService.CreateDeviceConfig(c, &device_config.CreateDeviceConfigOpts{
		ProductionOrderID:      req.ProductionOrderID,
		ProductionPlanID:       req.ProductionPlanID,
		DeviceID:               req.DeviceID,
		Color:                  req.Color,
		MaThongSoMay:           req.MaThongSoMay,
		MaTaiLieuHuongDan:      req.MaTaiLieuHuongDan,
		NgayHieuLuc:            req.NgayHieuLuc,
		StageID:                req.StageID,
		MaSanPham:              req.MaSanPham,
		MaSanPhamNoiBo:         req.MaSanPhamNoiBo,
		ChuKyIn:                req.ChuKyIn,
		ThoiGianChuanBi:        req.ThoiGianChuanBi,
		TenMauMuc:              req.TenMauMuc,
		TenLoaiMuc:             req.TenLoaiMuc,
		SoThuTuIn:              req.SoThuTuIn,
		MaPhim:                 req.MaPhim,
		MaMauMuc:               req.MaMauMuc,
		TinhTrangMuc:           req.TinhTrangMuc,
		DienTichPhuMuc:         req.DienTichPhuMuc,
		DungMoi:                req.DungMoi,
		NhietDoSay:             req.NhietDoSay,
		ThoiGianSay:            req.ThoiGianSay,
		GhiChu:                 req.GhiChu,
		MaKhung:                req.MaKhung,
		ThongSoLua:             req.ThongSoLua,
		KhoangCachKhungInBanIn: req.KhoangCachKhungInBanIn,
		CachIn:                 req.CachIn,
		CungDao:                req.CungDao,
		DoBenDao:               req.DoBenDao,
		DoNghiengDao:           req.DoNghiengDao,
		TocDoDao:               req.TocDoDao,
		TocDo:                  req.TocDo,
		Description:            req.Description,
		DeviceType:             req.DeviceType,
		DeviceConfig:           req.DeviceConfig,
		InkID:                  req.InkID,
		CreatedBy:              userId,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.CreateDeviceConfigResponse{
		ID: id,
	})
}

func (s deviceConfigController) EditDeviceConfig(c *gin.Context) {
	req := &dto.EditDeviceConfigRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.deviceConfigService.EditDeviceConfig(c, &device_config.EditDeviceConfigOpts{
		ID:           req.ID,
		DeviceConfig: req.DeviceConfig,
		Color:        req.Color,
		Description:  req.Description,
		//Search:                 req.Search,
		DeviceID:               req.DeviceID,
		MaThongSoMay:           req.MaThongSoMay,
		MaTaiLieuHuongDan:      req.MaTaiLieuHuongDan,
		NgayHieuLuc:            req.NgayHieuLuc,
		StageID:                req.StageID,
		MaSanPham:              req.MaSanPham,
		MaSanPhamNoiBo:         req.MaSanPhamNoiBo,
		ChuKyIn:                req.ChuKyIn,
		ThoiGianChuanBi:        req.ThoiGianChuanBi,
		TenMauMuc:              req.TenMauMuc,
		TenLoaiMuc:             req.TenLoaiMuc,
		SoThuTuIn:              req.SoThuTuIn,
		MaPhim:                 req.MaPhim,
		MaMauMuc:               req.MaMauMuc,
		TinhTrangMuc:           req.TinhTrangMuc,
		DienTichPhuMuc:         req.DienTichPhuMuc,
		DungMoi:                req.DungMoi,
		NhietDoSay:             req.NhietDoSay,
		ThoiGianSay:            req.ThoiGianSay,
		GhiChu:                 req.GhiChu,
		MaKhung:                req.MaKhung,
		ThongSoLua:             req.ThongSoLua,
		KhoangCachKhungInBanIn: req.KhoangCachKhungInBanIn,
		CachIn:                 req.CachIn,
		CungDao:                req.CungDao,
		DoBenDao:               req.DoBenDao,
		DoNghiengDao:           req.DoNghiengDao,
		TocDoDao:               req.TocDoDao,
		TocDo:                  req.TocDo,
		InkID:                  req.InkID,
		UpdatedBy:              interceptor.UserIDFromCtx(c),
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.EditDeviceConfigResponse{})
}

func (s deviceConfigController) DeleteDeviceConfig(c *gin.Context) {
	req := &dto.DeleteDeviceConfigRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.deviceConfigService.Delete(c, req.ID)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.DeleteDeviceConfigResponse{})
}

func (s deviceConfigController) FindDeviceConfigs(c *gin.Context) {
	req := &dto.FindDeviceConfigsRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	deviceConfigs, cnt, err := s.deviceConfigService.FindDeviceConfigs(c, &device_config.FindDeviceConfigsOpts{
		IDs:               req.Filter.IDs,
		Search:            req.Filter.Search,
		ProductionOrderID: req.Filter.ProductionOrderID,
		ProductionPlanID:  req.Filter.ProductionPlanID,
		DeviceType:        req.Filter.DeviceType,
	}, &repository.Sort{
		Order: repository.SortOrderDESC,
		By:    "ID",
	}, req.Paging.Limit, req.Paging.Offset)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	deviceConfigResp := make([]*dto.DeviceConfig, 0, len(deviceConfigs))
	for _, f := range deviceConfigs {
		deviceConfigResp = append(deviceConfigResp, toDeviceConfigResp(f))
	}

	transportutil.SendJSONResponse(c, &dto.FindDeviceConfigsResponse{
		DeviceConfigs: deviceConfigResp,
		Total:         cnt.Count,
	})
}

func toDeviceConfigResp(f *device_config.Data) *dto.DeviceConfig {
	return &dto.DeviceConfig{
		ID:                     f.ID,
		ProductionOrderID:      f.ProductionOrderID.String,
		ProductionOrderName:    f.ProductionOrderName.String,
		ProductionPlanID:       f.ProductionPlanID.String,
		DeviceID:               f.DeviceID.String,
		DeviceName:             f.DeviceName.String,
		DeviceType:             f.DeviceType,
		DeviceCode:             f.DeviceCode.String,
		Color:                  f.Color.String,
		MaThongSoMay:           f.MaThongSoMay.String,
		MaTaiLieuHuongDan:      f.MaTaiLieuHuongDan.String,
		NgayHieuLuc:            f.NgayHieuLuc.Time,
		StageID:                f.StageID.String,
		MaSanPham:              f.MaSanPham.String,
		MaSanPhamNoiBo:         f.MaSanPhamNoiBo.String,
		ChuKyIn:                f.ChuKyIn.String,
		ThoiGianChuanBi:        f.ThoiGianChuanBi.String,
		TenMauMuc:              f.TenMauMuc.String,
		TenLoaiMuc:             f.TenLoaiMuc.String,
		SoThuTuIn:              f.SoThuTuIn.String,
		MaPhim:                 f.MaPhim.String,
		MaMauMuc:               f.MaMauMuc.String,
		TinhTrangMuc:           f.TinhTrangMuc.String,
		DienTichPhuMuc:         f.DienTichPhuMuc.String,
		DungMoi:                f.DungMoi.String,
		NhietDoSay:             f.NhietDoSay.String,
		ThoiGianSay:            f.ThoiGianSay.String,
		GhiChu:                 f.GhiChu.String,
		MaKhung:                f.MaKhung.String,
		ThongSoLua:             f.ThongSoLua.String,
		KhoangCachKhungInBanIn: f.KhoangCachKhungInBanIn.String,
		CachIn:                 f.CachIn.String,
		CungDao:                f.CungDao.String,
		DoBenDao:               f.DoBenDao.String,
		DoNghiengDao:           f.DoNghiengDao.String,
		TocDoDao:               f.TocDoDao.String,
		TocDo:                  f.TocDo.String,
		Description:            f.Description.String,
		InkID:                  f.InkID.String,
		CreatedBy:              f.CreatedBy,
		CreatedAt:              f.CreatedAt,
		UpdatedAt:              f.UpdatedAt,
	}
}

func RegisterDeviceConfigController(
	r *gin.RouterGroup,
	deviceConfigService device_config.Service,
) {
	g := r.Group("device-config")

	var c DeviceConfigController = &deviceConfigController{
		deviceConfigService: deviceConfigService,
	}

	routeutil.AddEndpoint(
		g,
		"create",
		c.CreateDeviceConfig,
		&dto.CreateDeviceConfigRequest{},
		&dto.CreateDeviceConfigResponse{},
		"Create deviceConfig",
	)

	routeutil.AddEndpoint(
		g,
		"edit",
		c.EditDeviceConfig,
		&dto.EditDeviceConfigRequest{},
		&dto.EditDeviceConfigResponse{},
		"Edit deviceConfig",
	)

	routeutil.AddEndpoint(
		g,
		"delete",
		c.DeleteDeviceConfig,
		&dto.DeleteDeviceConfigRequest{},
		&dto.DeleteDeviceConfigResponse{},
		"delete deviceConfig",
	)

	routeutil.AddEndpoint(
		g,
		"find",
		c.FindDeviceConfigs,
		&dto.FindDeviceConfigsRequest{},
		&dto.FindDeviceConfigsResponse{},
		"Find deviceConfigs",
	)
}
