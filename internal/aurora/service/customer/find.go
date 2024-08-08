package customer

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

func (c *customerService) FindCustomers(ctx context.Context, opts *FindCustomersOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error) {
	filter := &repository.SearchCustomerOpts{
		IDs:    opts.IDs,
		Name:   opts.Name,
		Phone:  opts.Phone,
		Code:   opts.Code,
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	}
	customers, err := c.customerRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	total, err := c.customerRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*Data, 0, len(customers))
	for _, customer := range customers {
		results = append(results, &Data{
			CustomerData: customer,
		})
	}
	return results, total, nil
}

type FindCustomersOpts struct {
	IDs   []string
	Name  string
	Phone string
	Code  string
}
