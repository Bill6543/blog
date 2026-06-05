package validator

import (
	"errors"
	"regexp"
	"strings"
)

// ValidatePassword 验证密码强度
// 要求：
// 1. 长度 6-50 字符
// 2. 至少包含一个字母
// 3. 至少包含一个数字
// 4. 不能包含空格
func ValidatePassword(password string) error {
	// 1. 验证长度
	if len(password) < 6 || len(password) > 50 {
		return errors.New("密码长度必须在 6-50 个字符之间")
	}

	// 2. 不能包含空格
	if strings.Contains(password, " ") {
		return errors.New("密码不能包含空格")
	}

	// 3. 至少包含一个字母
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	if !hasLetter {
		return errors.New("密码必须包含至少一个字母")
	}

	// 4. 至少包含一个数字
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return errors.New("密码必须包含至少一个数字")
	}

	return nil
}

// ValidateUsername 验证用户名
func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 50 {
		return errors.New("用户名长度必须在 3-50 个字符之间")
	}
	return nil
}

// ValidateEmail 验证邮箱（补充 binding 的不足）
func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("邮箱不能为空")
	}
	// 可以添加更严格的验证
	return nil
}
