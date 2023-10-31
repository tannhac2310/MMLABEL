package passwordutil

import "golang.org/x/crypto/bcrypt"

// HashAndSalt ...
func HashAndSalt(pwd string) (string, error) {
	pwdByte := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(pwdByte, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePasswords ...
func ComparePasswords(hashedPwd string, pwd string) bool {
	plainPwd := []byte(pwd)
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	return err == nil
}
