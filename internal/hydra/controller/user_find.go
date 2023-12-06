package controller

import (
	"github.com/gin-gonic/gin"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"

	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/dto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/user"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

func (u *userController) FindUsers(c *gin.Context) {
	req := &dto.FindUsersRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	users, total, err := u.userService.SearchUsers(c, &user.SearchUsersOpts{
		IDs:         req.Filter.IDs,
		Name:        req.Filter.Name,
		Department:  req.Filter.Department,
		Departments: req.Filter.Departments,
		NotIDs:      req.Filter.NotIDs,
		NotRoleIDs:  req.Filter.NotRoleIDs,
		Search:      req.Filter.Search,
		Type:        req.Filter.Type,
		PhoneNumber: req.Filter.PhoneNumber,
		Email:       req.Filter.Email,
		GroupID:     req.Filter.GroupID,
		RoleID:      req.Filter.RoleID,
	}, req.Paging.Limit, req.Paging.Offset)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	userResp := make([]*dto.User, 0, len(users))
	for _, f := range users {
		userResp = append(userResp, toUserResp(f))
	}

	transportutil.SendJSONResponse(c, &dto.FindUsersResponse{
		Users: userResp,
		Total: total.Count,
	})
}

func toUserResp(e *repository.UserData) *dto.User {
	return &dto.User{
		ID:          e.ID,
		Name:        e.Name,
		Avatar:      e.Avatar,
		PhoneNumber: e.PhoneNumber,
		Email:       e.Email,
		Address:     e.Address,
		Type:        e.Type,
		Status:      e.Status,
		Code:        e.Code,
		Departments: e.Departments.String,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
