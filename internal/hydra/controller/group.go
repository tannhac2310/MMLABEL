package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/dto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/group"
)

type GroupController interface {
	CreateGroup(c *gin.Context)
	EditGroup(c *gin.Context)
	FindGroup(c *gin.Context)
	FindGroupByID(c *gin.Context)
	AddGroupsForUser(c *gin.Context)
	RemoveGroupsForUser(c *gin.Context)
}

type groupController struct {
	groupService group.Service
}

func RegisterGroupController(
	r *gin.RouterGroup,
	groupService group.Service,
) {
	g := r.Group("group")

	var c GroupController = &groupController{
		groupService: groupService,
	}

	routeutil.AddEndpoint(
		g,
		"create-group",
		c.CreateGroup,
		&dto.CreateGroupRequest{},
		&dto.CreateGroupResponse{},
		"Create group",
	)

	routeutil.AddEndpoint(
		g,
		"edit-group",
		c.EditGroup,
		&dto.EditGroupRequest{},
		&dto.EditGroupResponse{},
		"Edit group",
	)

	routeutil.AddEndpoint(
		g,
		"find-group",
		c.FindGroup,
		&dto.FindGroupRequest{},
		&dto.FindGroupResponse{},
		"Find group",
	)

	routeutil.AddEndpoint(
		g,
		"find-group-by-id",
		c.FindGroupByID,
		&dto.FindGroupByIDRequest{},
		&dto.FindGroupByIDResponse{},
		"Find one group",
	)

	routeutil.AddEndpoint(
		g,
		"add-groups-for-user",
		c.AddGroupsForUser,
		&dto.AddGroupsForUserRequest{},
		&dto.AddGroupsForUserResponse{},
		"Add groups for user",
	)

	routeutil.AddEndpoint(
		g,
		"remove-groups-for-user",
		c.RemoveGroupsForUser,
		&dto.RemoveGroupsForUserRequest{},
		&dto.RemoveGroupsForUserResponse{},
		"Remove groups for user",
	)
}
