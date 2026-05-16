package utils

import (
	"auth-micro-service/pkg/shortcut"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(
	`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`,
)

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)

	if !emailRegex.MatchString(email) {
		return shortcut.ErrValidateEmail
	}

	return nil
}
