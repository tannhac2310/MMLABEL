package dto

import "mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"

type (
	Group struct {
		ID    string   `json:"id"`
		Name  string   `json:"name"`
		Roles []string `json:"roles"`
	}

	// CreateGroup endpoint
	CreateGroupRequest struct {
		ID    string   `json:"id" binding:"required"`
		Name  string   `json:"name" binding:"required"`
		Roles []string `json:"roles"`
	}

	CreateGroupResponse struct {
		ID string `json:"id"`
	}

	// EditGroup endpoint
	EditGroupRequest struct {
		ID    string   `json:"id" binding:"required"`
		Name  string   `json:"name" binding:"required"`
		Roles []string `json:"roles"`
	}

	EditGroupResponse struct {
	}

	GroupFilter struct {
		IDs  []string `json:"ids"`
		Name string   `json:"name"`
	}

	// FindGroup endpoint
	FindGroupRequest struct {
		Filter *GroupFilter      `json:"filter"`
		Paging *commondto.Paging `json:"paging" binding:"required"`
	}

	FindGroupResponse struct {
		Groups   []*Group          `json:"groups"`
		NextPage *commondto.Paging `json:"nextPage"`
	}

	// FindOneGroup endpoint
	FindGroupByIDRequest struct {
		ID string `json:"id" binding:"required"`
	}

	FindGroupByIDResponse Group

	// AddGroupsForUser
	AddGroupsForUserRequest struct {
		UserID   string   `json:"userId" binding:"required"`
		GroupIDs []string `json:"groupIDs" binding:"required"`
	}

	AddGroupsForUserResponse struct{}

	// RemoveGroupsForUser
	RemoveGroupsForUserRequest struct {
		UserID   string   `json:"userId" binding:"required"`
		GroupIDs []string `json:"groupIDs" binding:"required"`
	}

	RemoveGroupsForUserResponse struct{}
)
