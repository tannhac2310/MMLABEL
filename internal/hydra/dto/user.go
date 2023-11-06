package dto

import (
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type (
	User struct {
		ID          string          `json:"id"`
		Name        string          `json:"name"`
		Avatar      string          `json:"avatar"`
		PhoneNumber string          `json:"phoneNumber"`
		Email       string          `json:"email"`
		Address     string          `json:"address"`
		Code        string          `json:"code"`
		Departments string          `json:"departments"`
		Type        enum.UserType   `json:"type"`
		Status      enum.UserStatus `json:"status"`
		CreatedAt   time.Time       `json:"createdAt"`
		UpdatedAt   time.Time       `json:"updatedAt"`
	}

	UserDetail struct {
		ID          string          `json:"id"`
		Name        string          `json:"name"`
		Avatar      string          `json:"avatar"`
		Address     string          `json:"address"`
		PhoneNumber string          `json:"phoneNumber"`
		Email       string          `json:"email"`
		Type        enum.UserType   `json:"type"`
		Status      enum.UserStatus `json:"status"`
		GroupIDs    []string        `json:"groupIds"`
		RoleIDs     []string        `json:"roleIds"`
		CreatedAt   time.Time       `json:"createdAt"`
		UpdatedAt   time.Time       `json:"updatedAt"`
	}

	GetCurrentProfileRequest struct {
	}

	GetCurrentProfileResponse UserProfile

	CheckUserNameRequest struct {
		UserName string `json:"userName" binding:"required"`
	}

	CheckUserNameResponse struct {
		Existed bool `json:"existed"`
	}

	RegisterLoginAccountRequest struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}

	RegisterLoginAccountResponse struct {
	}

	CreateUserRequest struct {
		Name        string        `json:"name" binding:"required"`
		Avatar      string        `json:"avatar"`
		Address     string        `json:"address"  binding:"required"`
		PhoneNumber string        `json:"phoneNumber" binding:"required"`
		Email       string        `json:"email" binding:"required,email"`
		Type        enum.UserType `json:"type" binding:"required"`
		Password    string        `json:"password"`
	}

	CreateUserResponse struct {
		ID string `json:"id"`
	}

	EditUserRequest struct {
		ID          string
		Name        string          `json:"name" binding:"required"`
		Avatar      string          `json:"avatar"`
		Status      enum.UserStatus `json:"status"  binding:"required"`
		Address     string          `json:"address"  binding:"required"`
		PhoneNumber string          `json:"phoneNumber"`
		Email       string          `json:"email" binding:"required,email"`
	}

	EditUserResponse struct{}

	ChangeUserPasswordRequest struct {
		UserID   string `json:"userId" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	ChangeUserPasswordResponse struct{}

	FindUserByIDRequest struct {
		ID string `json:"id" binding:"required"`
	}

	FindUserByIDResponse UserDetail

	UserFilter struct {
		IDs         []string      `json:"ids"`
		Name        string        `json:"name"`
		NotIDs      []string      `json:"notIds"`
		NotRoleIDs  []string      `json:"notRoleIDs"`
		Search      string        `json:"search"`
		Type        enum.UserType `json:"type"`
		PhoneNumber string        `json:"phoneNumber"`
		Email       string        `json:"email"`
		GroupID     string        `json:"groupId"`
		RoleID      string        `json:"roleId"`
	}

	FindUsersRequest struct {
		Filter *UserFilter       `json:"filter"`
		Paging *commondto.Paging `json:"paging" binding:"required"`
	}

	FindUsersResponse struct {
		Users []*User `json:"users"`
		Total int64   `json:"total"`
	}

	ChangePasswordRequest struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		Password    string `json:"password" binding:"required"`
	}

	ChangePasswordResponse struct{}

	UpdateProfileRequest struct {
		Name        string `json:"name" binding:"required"`
		Address     string `json:"address"`
		PhoneNumber string `json:"phoneNumber"`
		Avatar      string `json:"avatar"`
		Email       string `json:"email" binding:"required,email"`
	}

	UpdateProfileResponse struct{}

	UpdateFCMTokenRequest struct {
		DeviceID string `json:"deviceId"`
		Token    string `json:"token"`
	}

	UpdateFCMTokenResponse struct {
	}
)
type DeleteUserRequest struct {
	ID string `json:"id"`
}

type DeleteUserResponse struct{}
