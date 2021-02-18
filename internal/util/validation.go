package util

import (
	"fmt"
	"strings"
)

// ValidateUsernameRequirements validates username requirements
func ValidateUsernameRequirements(username string) (bool, string) {
	if strings.Contains(username, "@") {
		return false, "username is not to contain '@'"
	}
	return true, ""
}

// ValidateMinPasswordRequirements validates password requirements
func ValidateMinPasswordRequirements(password string) (bool, string) {
	minPasswordLength := 3
	if len(password) < minPasswordLength {
		return false, fmt.Sprintf("password is to be greater than %d", minPasswordLength)
	}
	// more to come!!!
	return true, ""
}

// ValidateEmailRequirements validates email requirements
func ValidateEmailRequirements(testEmail string) (bool, string) {
	split := strings.Split(testEmail, "@")
	if len(split) != 2 {
		return false, "email must have an @"
	}
	if !strings.Contains(split[1], ".") {
		return false, "email must have ."
	}
	return true, ""
}
