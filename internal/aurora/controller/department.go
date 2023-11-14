package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/dto"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/department"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/interceptor"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

type DepartmentController interface {
	CreateDepartment(c *gin.Context)
	EditDepartment(c *gin.Context)
	DeleteDepartment(c *gin.Context)
	FindDepartments(c *gin.Context)
}

type departmentController struct {
	departmentService department.Service
}

func (s departmentController) CreateDepartment(c *gin.Context) {
	req := &dto.CreateDepartmentRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	userID := interceptor.UserIDFromCtx(c)

	id, err := s.departmentService.CreateDepartment(c, &department.CreateDepartmentOpts{
		ParentID:  req.ParentID,
		Name:      req.Name,
		ShortName: req.ShortName,
		Code:      req.Code,
		Priority:  req.Priority,
		CreatedBy: userID,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.CreateDepartmentResponse{
		ID: id,
	})
}

func (s departmentController) EditDepartment(c *gin.Context) {
	req := &dto.EditDepartmentRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.departmentService.EditDepartment(c, &department.EditDepartmentOpts{
		ID:        req.ID,
		ParentID:  req.ParentID,
		Name:      req.Name,
		ShortName: req.ShortName,
		Code:      req.Code,
		Priority:  req.Priority,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.EditDepartmentResponse{})
}

func (s departmentController) DeleteDepartment(c *gin.Context) {
	req := &dto.DeleteDepartmentRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.departmentService.Delete(c, req.ID)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.DeleteDepartmentResponse{})
}

func (s departmentController) FindDepartments(c *gin.Context) {
	req := &dto.FindDepartmentsRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	departments, cnt, err := s.departmentService.FindDepartments(c, &department.FindDepartmentsOpts{
		Name: req.Filter.Name,
	}, &repository.Sort{
		Order: repository.SortOrderASC,
		By:    "ID",
	}, req.Paging.Limit, req.Paging.Offset)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	departmentResp := make([]*dto.Department, 0, len(departments))
	for _, f := range departments {
		departmentResp = append(departmentResp, toDepartmentResp(f))
	}

	transportutil.SendJSONResponse(c, &dto.FindDepartmentsResponse{
		Departments: departmentResp,
		Total:       cnt.Count,
	})
}

func toDepartmentResp(f *department.Data) *dto.Department {
	return &dto.Department{
		ID:        f.ID,
		ParentID:  f.ParentID.String,
		Name:      f.Name,
		ShortName: f.ShortName,
		Code:      f.Code,
		Priority:  f.Priority,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

func RegisterDepartmentController(
	r *gin.RouterGroup,
	departmentService department.Service,
) {
	g := r.Group("department")

	var c DepartmentController = &departmentController{
		departmentService: departmentService,
	}

	routeutil.AddEndpoint(
		g,
		"create",
		c.CreateDepartment,
		&dto.CreateDepartmentRequest{},
		&dto.CreateDepartmentResponse{},
		"Create department",
	)

	routeutil.AddEndpoint(
		g,
		"edit",
		c.EditDepartment,
		&dto.EditDepartmentRequest{},
		&dto.EditDepartmentResponse{},
		"Edit department",
	)

	routeutil.AddEndpoint(
		g,
		"delete",
		c.DeleteDepartment,
		&dto.DeleteDepartmentRequest{},
		&dto.DeleteDepartmentResponse{},
		"delete department",
	)

	routeutil.AddEndpoint(
		g,
		"find",
		c.FindDepartments,
		&dto.FindDepartmentsRequest{},
		&dto.FindDepartmentsResponse{},
		"Find departments",
	)
}
