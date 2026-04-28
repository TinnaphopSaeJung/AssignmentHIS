package utils

import "regexp"

var (
	hasLetter  = regexp.MustCompile(`[A-Za-z]`)
	hasNumber  = regexp.MustCompile(`\d`)
	hasSpecial = regexp.MustCompile(`[@$!%*#?&]`)
)

var thaiRegex = regexp.MustCompile(`^[ก-๙]+$`)

var engRegex = regexp.MustCompile(`^[A-Za-z-]+$`)

func IsValidPassword(pw string) bool {
	if len(pw) < 8 {
		return false
	}
	if !hasLetter.MatchString(pw) {
		return false
	}
	if !hasNumber.MatchString(pw) {
		return false
	}
	if !hasSpecial.MatchString(pw) {
		return false
	}
	return true
}

func IsThaiName(s string) bool {
	if s == "" {
		return true
	}
	return thaiRegex.MatchString(s)
}

func IsEnglishName(s string) bool {
	if s == "" {
		return true
	}
	return engRegex.MatchString(s)
}
