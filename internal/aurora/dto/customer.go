package dto

import (
	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type CustomerFilter struct {
	Name  string `json:"name"`
	Code  string `json:"code"`
	Phone string `json:"phone"`
}

type FindCustomersRequest struct {
	Filter *CustomerFilter   `json:"filter" binding:"required"`
	Paging *commondto.Paging `json:"paging" binding:"required"`
}

type FindCustomersResponse struct {
	Customers []*Customer `json:"customers"`
	Total     int64       `json:"total"`
}
type Customer struct {
	ID                 string              `json:"id,omitempty"`
	Name               string              `json:"name,omitempty"`
	Tax                string              `json:"tax,omitempty"`
	Code               string              `json:"code,omitempty"`
	Country            string              `json:"country,omitempty"`
	Province           string              `json:"province,omitempty"`
	Address            string              `json:"address,omitempty"`
	Fax                string              `json:"fax,omitempty"`
	CompanyWebsite     string              `json:"companyWebsite,omitempty"`
	CompanyPhone       string              `json:"companyPhone,omitempty"`
	ContactPersonName  string              `json:"contactPersonName,omitempty"`
	ContactPersonEmail string              `json:"contactPersonEmail,omitempty"`
	ContactPersonPhone string              `json:"contactPersonPhone,omitempty"`
	ContactPersonRole  string              `json:"contactPersonRole,omitempty"`
	Note               string              `json:"note,omitempty"`
	Status             enum.CustomerStatus `json:"status,omitempty"`
}

type CreateCustomerRequest struct {
	Name               string `json:"name,omitempty" binding:"required"`
	Tax                string `json:"tax,omitempty"`
	Code               string `json:"code,omitempty"`
	Country            string `json:"country,omitempty" `
	Province           string `json:"province,omitempty" `
	Address            string `json:"address,omitempty" `
	Fax                string `json:"fax,omitempty"`
	CompanyWebsite     string `json:"companyWebsite,omitempty"`
	CompanyPhone       string `json:"companyPhone,omitempty"`
	CompanyEmail       string `json:"companyEmail,omitempty"`
	ContactPersonName  string `json:"contactPersonName,omitempty" `
	ContactPersonEmail string `json:"contactPersonEmail,omitempty" `
	ContactPersonPhone string `json:"contactPersonPhone,omitempty" `
	ContactPersonRole  string `json:"contactPersonRole,omitempty"`
	Note               string `json:"note,omitempty"`
}

type CreateCustomerResponse struct {
	ID string `json:"id"`
}

type EditCustomerRequest struct {
	ID                 string `json:"id,omitempty" binding:"required"`
	Name               string `json:"name,omitempty" binding:"required"`
	Tax                string `json:"tax,omitempty"`
	Code               string `json:"code,omitempty"`
	Country            string `json:"country,omitempty" `
	Province           string `json:"province,omitempty" `
	Address            string `json:"address,omitempty" `
	Fax                string `json:"fax,omitempty"`
	CompanyWebsite     string `json:"companyWebsite,omitempty"`
	CompanyPhone       string `json:"companyPhone,omitempty"`
	CompanyEmail       string `json:"companyEmail,omitempty"`
	ContactPersonName  string `json:"contactPersonName,omitempty" `
	ContactPersonEmail string `json:"contactPersonEmail,omitempty" `
	ContactPersonPhone string `json:"contactPersonPhone,omitempty" `
	ContactPersonRole  string `json:"contactPersonRole,omitempty"`
	Note               string `json:"note,omitempty"`
}

type EditCustomerResponse struct {
}

type DeleteCustomerRequest struct {
	ID string `json:"id"`
}

type DeleteCustomerResponse struct {
}
