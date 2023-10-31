package dto

type UserProfile struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
}
