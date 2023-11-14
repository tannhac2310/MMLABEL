package department

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

func (c *departmentService) FindDepartments(ctx context.Context, opts *FindDepartmentsOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error) {
	filter := &repository.SearchDepartmentsOpts{
		Name:   opts.Name,
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	}
	departments, err := c.departmentRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	total, err := c.departmentRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*Data, 0, len(departments))
	for _, department := range departments {
		if err != nil {
			return nil, nil, err
		}
		results = append(results, &Data{
			DepartmentData: department,
		})
	}
	return results, total, nil
}

type FindDepartmentsOpts struct {
	Name string
}
