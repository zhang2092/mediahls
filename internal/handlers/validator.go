package handlers

import (
	"regexp"
	"strings"
)

var (
	rxPhone    = regexp.MustCompile(`^(13|14|15|16|17|18|19)\d{9}$`)
	rxEmail    = regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)
	rxUsername = regexp.MustCompile(`^[a-z0-9A-Z]{6,20}$`) // 6到20位（字母，数字）
	//rxPassword = regexp.MustCompile(`^(?=.*[a-zA-Z])(?=.*[0-9])[A-Za-z0-9]{8,18}$`) // 最少6位，包括至少1个大写字母，1个小写字母，1个数字，1个特殊字符
)

func ValidateRxPhone(phone string) bool {
	phone = strings.TrimSpace(phone)
	return rxPhone.MatchString(phone)
}

func ValidateRxEmail(email string) bool {
	email = strings.TrimSpace(email)
	return rxEmail.MatchString(email)
}

func ValidateRxUsername(username string) bool {
	username = strings.TrimSpace(username)
	return rxUsername.MatchString(username)
}

// func ValidateRxPassword(password string) bool {
// 	password = strings.TrimSpace(password)
// 	return rxPassword.MatchString(password)
// }

func ValidatePassword(password string) bool {
	password = strings.TrimSpace(password)
	return len(password) >= 8 && len(password) <= 20
}
