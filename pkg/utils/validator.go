package utils

import (
	"regexp"
	"time"
)

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

func IsValidDate(date string) bool {
	if date == "" {
		return true
	}
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}
