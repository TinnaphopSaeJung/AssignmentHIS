package utils

import "regexp"

var (
	hasLetter  = regexp.MustCompile(`[A-Za-z]`)
	hasNumber  = regexp.MustCompile(`\d`)
	hasSpecial = regexp.MustCompile(`[@$!%*#?&]`)
)

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
