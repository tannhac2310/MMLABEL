package dto

import (
	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
)

type CustomerFilter struct {
	Name  string `json:"name"`
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
	ID          string `json:"id"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Status      int16  `json:"status"`
	Type        int16  `json:"type"`
	Address     string `json:"address"`
}

type CreateCustomerRequest struct {
	Name        string `json:"name" binding:"required"`
	Avatar      string `json:"avatar"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Status      int16  `json:"status" binding:"required"`
	Type        int16  `json:"type"`
	Address     string `json:"address"`
}

type CreateCustomerResponse struct {
	ID string `json:"id"`
}

type EditCustomerRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name" binding:"required"`
	Avatar      string `json:"avatar"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Status      int16  `json:"status" binding:"required"`
	Type        int16  `json:"type"`
	Address     string `json:"address"`
}

type EditCustomerResponse struct {
}

type DeleteCustomerRequest struct {
	ID string `json:"id"`
}

type DeleteCustomerResponse struct {
}
