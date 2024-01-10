package device

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

func (c *deviceService) FindDevices(ctx context.Context, opts *FindDevicesOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error) {
	// get all children of parent department
	departments, err := c.stageRepo.Search(ctx, &repository.SearchStagesOpts{
		Limit:  1000,
		Offset: 0,
	})
	steps := make([]string, 0)
	if opts.Step != "" {
		steps = append(steps, opts.Step)
	}
	var step *repository.StageData

	for _, department := range departments {
		fmt.Print(department.ID, ",")
		if department.ID == opts.Step {
			step = department
		}
	}
	fmt.Println("step", step, opts.Step)
	parents := make([]*repository.StageData, 0)
	allDepartments := getAllParents(departments, step, parents)
	for _, department := range allDepartments {
		steps = append(steps, department.ID)
	}
	fmt.Println("steps", departments, steps)
	filter := &repository.SearchDevicesOpts{
		Name:   opts.Name,
		Steps:  steps,
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	}

	devices, err := c.deviceRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	total, err := c.deviceRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*Data, 0, len(devices))
	for _, device := range devices {
		if err != nil {
			return nil, nil, err
		}
		results = append(results, &Data{
			DeviceData: device,
		})
	}
	return results, total, nil
}

// parents := make([]*repository.DepartmentData, 0)
// getAllParents is recursive function to get all children of parent department by step (department)
func getAllParents(departments []*repository.StageData, step *repository.StageData, parents []*repository.StageData) []*repository.StageData {
	if step == nil {
		return parents
	}
	for _, department := range departments {
		if department.ID == step.ParentID.String {
			fmt.Println("department", step)
			parents = append(parents, department)
			parents = append(parents, getAllParents(departments, department, parents)...)
		}
	}
	return parents
}

type FindDevicesOpts struct {
	Name string
	Step string
}
