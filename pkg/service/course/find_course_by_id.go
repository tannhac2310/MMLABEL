package course

import (
	"context"
	"fmt"
)

func (c *courseService) FindCourseByID(ctx context.Context, id string) (*Course, error) {
	courses, total, err := c.FindCourses(ctx, &FindCoursesOpts{
		IDs: []string{id},
	}, 1, 0)

	if err != nil {
		return nil, err
	}
	if total.Count != 1 {
		return nil, fmt.Errorf("stage.Search:FindStageByID not found")
	}

	return courses[0], nil
}
