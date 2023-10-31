package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/dto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/auth"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/user"
)

type AuthController interface {
	LoginUserNamePassword(c *gin.Context)
	LoginFirebase(c *gin.Context)
	RefreshToken(c *gin.Context)
	RequestOTP(c *gin.Context)
	LoginOTP(c *gin.Context)
	ResetPassword(c *gin.Context)
}

type authController struct {
	authService auth.Service
	userService user.Service
}

func RegisterAuthController(
	r *gin.RouterGroup,
	authService auth.Service,
	userService user.Service,
) {
	g := r.Group("auth")

	var c AuthController = &authController{
		authService: authService,
		userService: userService,
	}

	routeutil.AddEndpoint(
		g,
		"login-username-password",
		c.LoginUserNamePassword,
		&dto.LoginUserNamePasswordRequest{},
		&dto.LoginUserNamePasswordResponse{},
		"Login with email and password",
		routeutil.RegisterOptionSkipAuth,
	)

	routeutil.AddEndpoint(
		g,
		"login-firebase",
		c.LoginFirebase,
		&dto.LoginFirebaseRequest{},
		&dto.LoginFirebaseResponse{},
		"Login with firebase",
		routeutil.RegisterOptionSkipAuth,
	)

	routeutil.AddEndpoint(
		g,
		"refresh-token",
		c.RefreshToken,
		&dto.RefreshTokenRequest{},
		&dto.RefreshTokenResponse{},
		"Refresh token",
		routeutil.RegisterOptionSkipAuth,
	)

	routeutil.AddEndpoint(
		g,
		"request-otp",
		c.RequestOTP,
		&dto.RequestOTPRequest{},
		&dto.RequestOTPResponse{},
		"will limit 5 sms per hours",
		routeutil.RegisterOptionSkipAuth,
	)

	routeutil.AddEndpoint(
		g,
		"login-otp",
		c.LoginOTP,
		&dto.LoginOTPRequest{},
		&dto.LoginOTPResponse{},
		"login with otp code",
		routeutil.RegisterOptionSkipAuth,
	)

	routeutil.AddEndpoint(
		g,
		"reset-password",
		c.ResetPassword,
		&dto.ResetPasswordRequest{},
		&dto.ResetPasswordResponse{},
		"login with otp code",
		routeutil.RegisterOptionSkipAuth,
	)
}
