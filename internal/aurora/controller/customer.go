package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/dto"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/customer"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/interceptor"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

type CustomerController interface {
	CreateCustomer(c *gin.Context)
	EditCustomer(c *gin.Context)
	DeleteCustomer(c *gin.Context)
	FindCustomers(c *gin.Context)
}

type customerController struct {
	customerService customer.Service
}

func (s customerController) CreateCustomer(c *gin.Context) {
	req := &dto.CreateCustomerRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	userID := interceptor.UserIDFromCtx(c)

	id, err := s.customerService.CreateCustomer(c, &customer.CreateCustomerOpts{
		Name:               req.Name,
		Tax:                req.Tax,
		Code:               req.Code,
		Country:            req.Country,
		Province:           req.Province,
		Address:            req.Address,
		Fax:                req.Fax,
		CompanyWebsite:     req.CompanyWebsite,
		CompanyPhone:       req.CompanyPhone,
		CompanyEmail:       req.CompanyEmail,
		ContactPersonName:  req.ContactPersonName,
		ContactPersonEmail: req.ContactPersonEmail,
		ContactPersonPhone: req.ContactPersonPhone,
		ContactPersonRole:  req.ContactPersonRole,
		Note:               req.Note,
		Data:               req.UserField,
		Status:             enum.CustomerStatusActivate,
		CreatedBy:          userID,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.CreateCustomerResponse{
		ID: id,
	})
}

func (s customerController) EditCustomer(c *gin.Context) {
	req := &dto.EditCustomerRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	userID := interceptor.UserIDFromCtx(c)

	err = s.customerService.EditCustomer(c, &customer.EditCustomerOpts{
		ID:                 req.ID,
		Name:               req.Name,
		Tax:                req.Tax,
		Code:               req.Code,
		Country:            req.Country,
		Province:           req.Province,
		Address:            req.Address,
		Fax:                req.Fax,
		CompanyWebsite:     req.CompanyWebsite,
		CompanyPhone:       req.CompanyPhone,
		Data:               req.UserField,
		ContactPersonName:  req.ContactPersonName,
		ContactPersonEmail: req.ContactPersonEmail,
		ContactPersonPhone: req.ContactPersonPhone,
		ContactPersonRole:  req.ContactPersonRole,
		Note:               req.Note,
		Status:             enum.CustomerStatusActivate,
		CreatedBy:          userID,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.EditCustomerResponse{})
}

func (s customerController) DeleteCustomer(c *gin.Context) {
	req := &dto.DeleteCustomerRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.customerService.Delete(c, req.ID)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.DeleteCustomerResponse{})
}

func (s customerController) FindCustomers(c *gin.Context) {
	req := &dto.FindCustomersRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	customers, cnt, err := s.customerService.FindCustomers(c, &customer.FindCustomersOpts{
		Name: req.Filter.Name,
		Code: req.Filter.Code,
		ID:   req.Filter.ID,
		IDs:  req.Filter.IDs,
	}, &repository.Sort{
		Order: repository.SortOrderDESC,
		By:    "ID",
	}, req.Paging.Limit, req.Paging.Offset)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	customerResp := make([]*dto.Customer, 0, len(customers))
	for _, f := range customers {
		customerResp = append(customerResp, toCustomerResp(f))
	}

	transportutil.SendJSONResponse(c, &dto.FindCustomersResponse{
		Customers: customerResp,
		Total:     cnt.Count,
	})
}

func toCustomerResp(f *customer.Data) *dto.Customer {
	return &dto.Customer{
		ID:                 f.ID,
		Name:               f.Name,
		Tax:                f.Tax.String,
		Code:               f.Code,
		Country:            f.Country,
		Province:           f.Province,
		Address:            f.Address,
		Fax:                f.Fax.String,
		CompanyWebsite:     f.CompanyWebsite.String,
		CompanyPhone:       f.CompanyPhone.String,
		ContactPersonName:  f.ContactPersonName.String,
		ContactPersonPhone: f.ContactPersonPhone.String,
		ContactPersonEmail: f.ContactPersonEmail.String,
		ContactPersonRole:  f.ContactPersonRole.String,
		CompanyEmail:       f.CompanyEmail.String,
		Note:               f.Note.String,
		UserField:          f.Data,
		Status:             f.Status,
	}
}

func RegisterCustomerController(
	r *gin.RouterGroup,
	customerService customer.Service,
) {
	g := r.Group("customer")

	var c CustomerController = &customerController{
		customerService: customerService,
	}

	routeutil.AddEndpoint(
		g,
		"create",
		c.CreateCustomer,
		&dto.CreateCustomerRequest{},
		&dto.CreateCustomerResponse{},
		"Create customer",
	)

	routeutil.AddEndpoint(
		g,
		"edit",
		c.EditCustomer,
		&dto.EditCustomerRequest{},
		&dto.EditCustomerResponse{},
		"Edit customer",
	)

	routeutil.AddEndpoint(
		g,
		"delete",
		c.DeleteCustomer,
		&dto.DeleteCustomerRequest{},
		&dto.DeleteCustomerResponse{},
		"delete customer",
	)

	routeutil.AddEndpoint(
		g,
		"find",
		c.FindCustomers,
		&dto.FindCustomersRequest{},
		&dto.FindCustomersResponse{},
		"Find customers",
	)
}
