package product

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	repository2 "mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

func (p productService) FindProduct(ctx context.Context, opts *FindProductOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error) {
	filter := repository.SearchProductOpts{
		IDs:           opts.IDs,
		Name:          opts.Name,
		CustomerID:    opts.CustomerID,
		SaleID:        opts.SaleID,
		ProductPlanID: opts.ProductPlanID,
		Limit:         limit,
		Offset:        offset,
		Sort:          nil,
	}

	// 1. Find products
	products, err := p.productRepo.Search(ctx, &filter)
	if err != nil {
		return nil, nil, fmt.Errorf("find products failed: %w", err)
	}
	// 2. Count products
	count, err := p.productRepo.Count(ctx, &filter)

	if err != nil {
		return nil, nil, fmt.Errorf("count products failed: %w", err)
	}
	customerIDs := make([]string, 0)
	userID := make([]string, 0)
	productionPlanIDs := make([]string, 0)
	// 3. Collect customerIDs and userIDs
	for _, product := range products {
		customerIDs = append(customerIDs, product.CustomerID)
		userID = append(userID, product.CreatedBy)
		userID = append(userID, product.UpdatedBy)
		if product.ProductionPlanID.String != "" {
			productionPlanIDs = append(productionPlanIDs, product.ProductionPlanID.String)
		}
	}

	// 4. Find customer data
	customerData, err := p.customerRepo.Search(ctx, &repository.SearchCustomerOpts{
		IDs:    customerIDs,
		Offset: 0,
		Limit:  int64(len(customerIDs)),
	})

	if err != nil {
		return nil, nil, fmt.Errorf("find customer data failed: %w", err)
	}

	// 5. Find user data
	userData, err := p.userRepo.Search(ctx, &repository2.SearchUsersOpts{
		IDs:    userID,
		Offset: 0,
		Limit:  int64(len(userID)),
	})

	if err != nil {
		return nil, nil, fmt.Errorf("find user data failed: %w", err)
	}

	// 6. Find production plan data
	productionPlanData, err := p.productionPlanRepo.Search(ctx, &repository.SearchProductionPlanOpts{
		IDs:    productionPlanIDs,
		Offset: 0,
		Limit:  int64(len(productionPlanIDs)),
	})

	if err != nil {
		return nil, nil, fmt.Errorf("find production plan data failed: %w", err)
	}

	// 7. Map customer data
	customerDataMap := make(map[string]*repository.CustomerData)
	for _, customer := range customerData {
		customerDataMap[customer.ID] = customer
	}

	// 8. Map user data
	userDataMap := make(map[string]*repository2.UserData)
	for _, user := range userData {
		userDataMap[user.ID] = user
	}

	// 9. Map production plan data
	productionPlanDataMap := make(map[string]*repository.ProductionPlanData)
	for _, productionPlan := range productionPlanData {
		productionPlanDataMap[productionPlan.ID] = productionPlan
	}

	// 10. Map product data
	var data []*Data
	for _, product := range products {
		// 10.1 find custom field value
		customFieldData, err := p.userFieldRepo.Search(ctx, &repository.SearchCustomFieldsOpts{
			EntityType: enum.CustomFieldTypeProduct,
			EntityId:   product.ID,
			Limit:      1000,
			Offset:     0,
		})
		if err != nil {
			return nil, nil, err
		}

		userFields := make([]*repository.CustomFieldData, 0)
		for _, datum := range customFieldData {
			userFields = append(userFields, datum)
		}

		data = append(data, &Data{
			ID:                 product.ID,
			Name:               product.Name,
			Code:               product.Code,
			CustomerID:         product.CustomerID,
			CustomerData:       customerDataMap[product.CustomerID],
			SaleID:             product.SaleID,
			Description:        product.Description,
			Data:               product.Data,
			UserFields:         userFields,
			CreatedAt:          product.CreatedAt,
			UpdatedAt:          product.UpdatedAt,
			CreatedBy:          product.CreatedBy,
			CreatedByName:      userDataMap[product.CreatedBy].Name,
			UpdatedBy:          product.UpdatedBy,
			UpdatedByName:      userDataMap[product.UpdatedBy].Name,
			ProductionPlanID:   product.ProductionPlanID.String,
			ProductionPlanData: productionPlanDataMap[product.ProductionPlanID.String],
		})
	}

	return data, count, nil
}
