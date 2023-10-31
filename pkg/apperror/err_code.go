//DO NOT EDIT: code generated from 'tools/gen-err-code.go'
package apperror

import "strconv"

type ErrCode int32

const (
	Langen_US = "en_US"
	Langvi_VN = "vi_VN"
)

var (
	ErrUnknown                  = &AppError{Code: 1000}
	ErrInvalidArgument          = &AppError{Code: 1001}
	ErrWrongEmail               = &AppError{Code: 1002}
	ErrWrongPassword            = &AppError{Code: 1003}
	ErrAccountInActive          = &AppError{Code: 1004}
	ErrUserNotFound             = &AppError{Code: 1005}
	ErrPermissionDenied         = &AppError{Code: 1006}
	ErrAccountExisted           = &AppError{Code: 1007}
	ErrInvalidOldPassword       = &AppError{Code: 1008}
	ErrInvalidOTP               = &AppError{Code: 1009}
	ErrUnauthenticated          = &AppError{Code: 2001}
	ErrAttendanceUserNotInClass = &AppError{Code: 3001}
	ErrNotFoundOrNotPermission  = &AppError{Code: 4001}
	ErrZaloAppIDExisted         = &AppError{Code: 4002}
)

var ErrCode_name = map[ErrCode]string{
	1000: "ErrUnknown",
	1001: "ErrInvalidArgument",
	1002: "ErrWrongEmail",
	1003: "ErrWrongPassword",
	1004: "ErrAccountInActive",
	1005: "ErrUserNotFound",
	1006: "ErrPermissionDenied",
	1007: "ErrAccountExisted",
	1008: "ErrInvalidOldPassword",
	1009: "ErrInvalidOTP",
	2001: "ErrUnauthenticated",
	3001: "ErrAttendanceUserNotInClass",
	4001: "ErrNotFoundOrNotPermission",
	4002: "ErrZaloAppIDExisted",
}

var ErrCode_value = map[string]ErrCode{
	"ErrUnknown":                  1000,
	"ErrInvalidArgument":          1001,
	"ErrWrongEmail":               1002,
	"ErrWrongPassword":            1003,
	"ErrAccountInActive":          1004,
	"ErrUserNotFound":             1005,
	"ErrPermissionDenied":         1006,
	"ErrAccountExisted":           1007,
	"ErrInvalidOldPassword":       1008,
	"ErrInvalidOTP":               1009,
	"ErrUnauthenticated":          2001,
	"ErrAttendanceUserNotInClass": 3001,
	"ErrNotFoundOrNotPermission":  4001,
	"ErrZaloAppIDExisted":         4002,
}

func (x ErrCode) String() string {
	s, ok := ErrCode_name[x]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}

var ErrCode_en_US = map[ErrCode]string{
	1000: "Something went wrong",
	1001: "Invalid argument",
	1002: "Wrong email",
	1003: "Wrong password",
	1004: "Your account is inactive",
	1005: "User not found",
	1006: "Permission denied",
	1007: "Account already existed",
	1008: "Old password does not match",
	1009: "Invalid OTP",
	2001: "Unauthenticated",
	3001: "Students who have not yet registered for the class",
	4001: "Data not found or you can not access it",
	4002: "Zalo app already existed",
}

var ErrCode_vi_VN = map[ErrCode]string{
	1000: "Lỗi không xác định",
	1001: "Dữ liệu đầu vào không hợp lệ",
	1002: "Sai email",
	1003: "Sai mật khẩu",
	1004: "Tài khoản chưa kích hoạt",
	1005: "Người dùng không tồn tại",
	1006: "Không có quyền sử dụng",
	1007: "Tài khoản đã tồn tại",
	1008: "Mật khẩu cũ không đúng",
	1009: "Mã OTP không hợp lệ",
	2001: "Lỗi xác thực",
	3001: " Học viên chưa đăng lý lớp học",
	4001: " Dữ liệu không tim thấy hoặc bạn không có quyền truy cập",
	4002: " Zalo app đã được kết nối",
}

func (x ErrCode) GetMessage(lang string) string {
	switch lang {
	case Langen_US:
		s, ok := ErrCode_en_US[x]
		if ok {
			return s
		}
	case Langvi_VN:
		s, ok := ErrCode_vi_VN[x]
		if ok {
			return s
		}
	default:
		s, ok := ErrCode_en_US[x]
		if ok {
			return s
		}
	}

	return strconv.Itoa(int(x))
}
