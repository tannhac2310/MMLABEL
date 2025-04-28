package production_plan

import (
	"context"
	"encoding/json"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type Analysis struct {
	Status enum.ProductionOrderStatus `json:"status"`
	Count  int64                      `json:"count"`
}

type FindProductionPlansOpts struct {
	IDs []string
	//CustomerID  string
	Search      string
	Name        string
	ProductName string
	ProductCode string
	Statuses    []enum.ProductionPlanStatus
	//UserID      string // todo add later
	Stage enum.ProductionPlanStage
}

func (c *productionPlanService) FindProductionPlans(ctx context.Context, opts *FindProductionPlansOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error) {
	fmt.Printf("opts: %+v\n", opts)
	filter := &repository.SearchProductionPlanOpts{
		IDs: opts.IDs,
		//CustomerID:  opts.CustomerID,
		Search:      opts.Search,
		Name:        opts.Name,
		ProductName: opts.ProductName,
		ProductCode: opts.ProductCode,
		Statuses:    opts.Statuses,
		//UserID:      opts.UserID,
		Stage:  opts.Stage,
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	}
	productionPlans, err := c.productionPlanRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	total, err := c.productionPlanRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	customerIds := make([]string, 0, len(productionPlans))

	for _, productionPlan := range productionPlans {
		// todo
		fmt.Println(productionPlan)
		//customerIds = append(customerIds, productionPlan.CustomerID)
	}
	// find customer name
	customerData, err := c.customerRepo.Search(ctx, &repository.SearchCustomerOpts{
		IDs:    customerIds,
		Limit:  int64(len(customerIds)),
		Offset: 0,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("c.customerRepo.FindByIDs: %w", err)
	}

	customerMap := make(map[string]*repository.CustomerData)
	for _, customer := range customerData {
		customerMap[customer.ID] = customer
	}
	results := make([]*Data, 0, len(productionPlans))
	for _, productionPlan := range productionPlans {
		// find custom field value
		customFieldData, err := c.customFieldRepo.Search(ctx, &repository.SearchCustomFieldsOpts{
			EntityType: enum.CustomFieldTypeProductionPlan,
			EntityId:   productionPlan.ID,
			Limit:      1000,
			Offset:     0,
		})
		if err != nil {
			return nil, nil, err
		}

		customFieldMap := make(map[string]string)
		userFields := make(map[string][]*repository.CustomFieldData)
		for _, datum := range customFieldData {
			customFieldMap[datum.Field] = datum.Value
			userFields[datum.Field] = append(userFields[datum.Field], datum)
		}

		// Handle san_pham_nguon field
		if productSource, ok := customFieldMap["san_pham_nguon"]; ok {
			// First try to parse as JSON array
			var productIDs []string
			if err := json.Unmarshal([]byte(productSource), &productIDs); err != nil {
				// If not a JSON array, try to parse as a single string
				productIDs = []string{productSource}
			}

			if len(productIDs) > 0 {
				// Find products by IDs
				products, err := c.productRepo.Search(ctx, &repository.SearchProductOpts{
					IDs:    productIDs,
					Limit:  int64(len(productIDs)),
					Offset: 0,
				})
				if err == nil && len(products) > 0 {
					// Collect all product codes and names
					productCodes := make([]string, 0, len(products))
					productNames := make([]string, 0, len(products))
					for _, product := range products {
						productCodes = append(productCodes, product.Code)
						productNames = append(productNames, product.Name)
					}

					// Convert to JSON string
					productCodesJSON, _ := json.Marshal(productCodes)
					productNamesJSON, _ := json.Marshal(productNames)

					// Update custom fields
					customFieldMap["product_code"] = string(productCodesJSON)
					customFieldMap["sale_survey_customer_product_code"] = string(productCodesJSON)
					customFieldMap["sale_survey_customer_product_name"] = string(productNamesJSON)
				}
			}
		}

		//_customerData := customerMap[productionPlan.CustomerID]
		results = append(results, &Data{
			ProductionPlanData: productionPlan,
			CustomData:         customFieldMap,
			UserFields:         userFields,
			//CustomerData:       _customerData,
		})
	}

	return results, total, nil
}

func (c *productionPlanService) FindProductionPlansWithNoPermission(ctx context.Context, opts *FindProductionPlansOpts, sort *repository.Sort, limit, offset int64) ([]*DataWithNoPermission, *repository.CountResult, error) {
	filter := &repository.SearchProductionPlanOpts{
		IDs: opts.IDs,
		//CustomerID:  opts.CustomerID,
		Name:        opts.Name,
		ProductName: opts.ProductName,
		ProductCode: opts.ProductCode,
		Statuses:    opts.Statuses,
		//UserID:      opts.UserID,
		Stage:  opts.Stage,
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	}
	productionPlans, err := c.productionPlanRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	total, err := c.productionPlanRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*DataWithNoPermission, 0, len(productionPlans))
	customerIds := make([]string, 0, len(productionPlans))

	for _, productionPlan := range productionPlans {
		// TODO
		fmt.Println(productionPlan)
		//customerIds = append(customerIds, productionPlan.CustomerID)
	}

	// find customer name
	customerData, err := c.customerRepo.Search(ctx, &repository.SearchCustomerOpts{
		IDs:    customerIds,
		Limit:  int64(len(customerIds)),
		Offset: 0,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("c.customerRepo.FindByIDs: %w", err)
	}

	customerMap := make(map[string]*repository.CustomerData)
	for _, customer := range customerData {
		customerMap[customer.ID] = customer
	}

	for _, productionPlan := range productionPlans {
		//_customerData := customerMap[productionPlan.CustomerID]
		results = append(results, &DataWithNoPermission{
			ProductionPlanData: productionPlan,
			//CustomerData:       _customerData,
		})
	}

	return results, total, nil
}
