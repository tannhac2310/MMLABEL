package dto

import "mmlabel.gitlab.com/mm-printing-backend/pkg/enum"

type UserProfile struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Code        string        `json:"code"`
	Departments string        `json:"departments"`
	Avatar      string        `json:"avatar"`
	PhoneNumber string        `json:"phoneNumber"`
	Email       string        `json:"email"`
	Type        enum.UserType `json:"type"`
}
