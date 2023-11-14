package dto

import (
	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"time"
)

type DepartmentFilter struct {
	Name string `json:"name"`
}

type FindDepartmentsRequest struct {
	Filter *DepartmentFilter `json:"filter" binding:"required"`
	Paging *commondto.Paging `json:"paging" binding:"required"`
}

type FindDepartmentsResponse struct {
	Departments []*Department `json:"departments"`
	Total       int64         `json:"total"`
}
type Department struct {
	ID        string    `json:"id"`
	ParentID  string    `json:"parentID"`
	Name      string    `json:"name"`
	ShortName string    `json:"shortName"`
	Code      string    `json:"code"`
	Priority  int64     `json:"priority"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateDepartmentRequest struct {
	ParentID  string `json:"parentID"`
	Name      string `json:"name" binding:"required"`
	ShortName string `json:"shortName" binding:"required"`
	Code      string `json:"code" binding:"required"`
	Priority  int64  `json:"priority"`
}

type CreateDepartmentResponse struct {
	ID string `json:"id"`
}

type EditDepartmentRequest struct {
	ID        string `json:"id" binding:"required"`
	ParentID  string `json:"parentID"`
	Name      string `json:"name" binding:"required"`
	ShortName string `json:"shortName" binding:"required"`
	Code      string `json:"code" binding:"required"`
	Priority  int64  `json:"priority"`
}

type EditDepartmentResponse struct {
}

type DeleteDepartmentRequest struct {
	ID string `json:"id"`
}

type DeleteDepartmentResponse struct {
}
