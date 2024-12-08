package dto

import (
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
)

type Product struct {
	ID                 string          `json:"id"`
	Name               string          `json:"name"`
	Code               string          `json:"code"`
	CustomerID         string          `json:"customerID"`
	CustomerData       *Customer       `json:"customerData"`
	ProductionPlanData *ProductionPlan `json:"productionPlanData"`
	ProductionPlanID   string          `json:"productionPlanID"`
	UserField          []*UserField    `json:"userField"`
	SaleID             string          `json:"saleID"`
	Description        string          `json:"description"`
	Data               any             `json:"data"`
	CreatedAt          time.Time       `json:"createdAt"`
	UpdatedAt          time.Time       `json:"updatedAt"`
	CreatedBy          string          `json:"createdBy"`
	CreatedByName      string          `json:"createdByName"`
	UpdatedBy          string          `json:"updatedBy"`
	UpdatedByName      string          `json:"updatedByName"`
}

type CreateProductRequest struct {
	Name        string       `json:"name"`
	Code        string       `json:"code"`
	CustomerID  string       `json:"customerID"`
	SaleID      string       `json:"saleID"`
	Description string       `json:"description"`
	Data        any          `json:"data"`
	UserField   []*UserField `json:"userField"`
}

type UserField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type CreateProductResponse struct {
	ID string `json:"id"`
}

type UpdateProductOpts struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Code        string       `json:"code"`
	CustomerID  string       `json:"customerID"`
	SaleID      string       `json:"saleID"`
	Description string       `json:"description"`
	Data        any          `json:"data"`
	UserField   []*UserField `json:"userField"`
}

type UpdateProductResponse struct {
}

type FindProductOpts struct {
	IDs                           []string `json:"ids"`
	Name                          string   `json:"name"`
	Code                          string   `json:"code"`
	CustomerID                    string   `json:"customerID"`
	SaleID                        string   `json:"saleID"`
	ProductionPlanID              string   `json:"productionPlanID"`
	ProductionOrderID             string   `json:"productionOrderID"`
	SaleSurveyCustomerProductName string   `json:"sale_survey_customer_product_name"`
	ProductName                   string   `json:"productName"`
	ProductCode                   string   `json:"productCode"`
	SaleSurveyCustomerProductCode string   `json:"sale_survey_customer_product_code"`
	SaleSurveyBravoCode           string   `json:"sale_survey_bravo_code"`
}
type SearchProductFilter struct {
	Filter FindProductOpts   `json:"filter" binding:"required"`
	Paging *commondto.Paging `json:"paging" binding:"required"`
}

type SearchProductResponse struct {
	Products []*Product `json:"products"`
	Total    int64      `json:"total"`
}

type DeleteProductRequest struct {
	ID string `json:"id"`
}

type DeleteProductResponse struct {
}
