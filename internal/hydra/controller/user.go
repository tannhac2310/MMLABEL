package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/dto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/group"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/role"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/user"
)

type UserController interface {
	GetCurrentProfile(c *gin.Context)
	CheckUserName(c *gin.Context)
	RegisterLoginAccount(c *gin.Context)
	CreateUser(c *gin.Context)
	EditUser(c *gin.Context)
	ChangeUserPassword(c *gin.Context)
	FindUserByID(c *gin.Context)
	FindUsers(c *gin.Context)
	ChangePassword(c *gin.Context)
	UpdateProfile(c *gin.Context)
	UpdateFCMToken(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type userController struct {
	userService  user.Service
	groupService group.Service
	roleService  role.Service
}

func RegisterUserController(
	r *gin.RouterGroup,
	userService user.Service,
	groupService group.Service,
	roleService role.Service,
) {
	g := r.Group("user")

	var c UserController = &userController{
		userService:  userService,
		groupService: groupService,
		roleService:  roleService,
	}

	routeutil.AddEndpoint(
		g,
		"get-current-profile",
		c.GetCurrentProfile,
		&dto.GetCurrentProfileRequest{},
		&dto.GetCurrentProfileResponse{},
		"Get current profile",
	)

	routeutil.AddEndpoint(
		g,
		"check-user-name",
		c.CheckUserName,
		&dto.CheckUserNameRequest{},
		&dto.CheckUserNameResponse{},
		"Check user name (email/phone_number) is existed or not",
		routeutil.RegisterOptionSkipAuth,
	)

	routeutil.AddEndpoint(
		g,
		"register-login-account",
		c.RegisterLoginAccount,
		&dto.RegisterLoginAccountRequest{},
		&dto.RegisterLoginAccountResponse{},
		"Register login account with userName and password",
	)

	routeutil.AddEndpoint(
		g,
		"create-user",
		c.CreateUser,
		&dto.CreateUserRequest{},
		&dto.CreateUserResponse{},
		"create user in cms",
	)

	routeutil.AddEndpoint(
		g,
		"edit-user",
		c.EditUser,
		&dto.EditUserRequest{},
		&dto.EditUserRequest{},
		"edit user in cms",
	)

	routeutil.AddEndpoint(
		g,
		"change-user-password",
		c.ChangeUserPassword,
		&dto.ChangeUserPasswordRequest{},
		&dto.ChangeUserPasswordResponse{},
		"change user password in cms",
	)

	routeutil.AddEndpoint(
		g,
		"find-user-by-id",
		c.FindUserByID,
		&dto.FindUserByIDRequest{},
		&dto.FindUserByIDResponse{},
		"find user by id",
	)

	routeutil.AddEndpoint(
		g,
		"find-users",
		c.FindUsers,
		&dto.FindUsersRequest{},
		&dto.FindUsersResponse{},
		"find users",
	)

	routeutil.AddEndpoint(
		g,
		"change-password",
		c.ChangePassword,
		&dto.ChangePasswordRequest{},
		&dto.ChangePasswordResponse{},
		"change password of current user",
	)

	routeutil.AddEndpoint(
		g,
		"update-profile",
		c.UpdateProfile,
		&dto.UpdateProfileRequest{},
		&dto.UpdateProfileResponse{},
		"update profile of current user",
	)

	routeutil.AddEndpoint(
		g,
		"update-fcm-token",
		c.UpdateFCMToken,
		&dto.UpdateFCMTokenRequest{},
		&dto.UpdateFCMTokenResponse{},
		"update fcm token",
	)

	routeutil.AddEndpoint(
		g,
		"delete-user",
		c.DeleteUser,
		&dto.DeleteUserRequest{},
		&dto.DeleteUserResponse{},
		"delete user",
	)
}
