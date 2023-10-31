package dto

type LoginResponse struct {
	Token        string       `json:"token"`
	RefreshToken string       `json:"refreshToken"`
	Profile      *UserProfile `json:"profile"`
	ACL          []string     `json:"acl"`
	Role         string       `json:"role"`
}

type LoginUserNamePasswordRequest struct {
	UserName string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserNamePasswordResponse LoginResponse

type LoginFirebaseRequest struct {
	IDToken string `json:"idToken"`
}

type LoginFirebaseResponse LoginResponse

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type RefreshTokenResponse LoginResponse

type RequestOTPRequest struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
}

type RequestOTPResponse struct{}

type LoginOTPRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	OTP         string `json:"otp" binding:"required"`
}

type LoginOTPResponse LoginResponse

type ResetPasswordRequest struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	OTP         string `json:"otp" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

type ResetPasswordResponse LoginResponse

type DeleteBannerRequest struct {
	ID string `json:"id"`
}
type DeleteBannerResponse struct{}
