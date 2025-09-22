package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidFullname = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidDataString(value string, minLen int, maxLen int) error {
	n := len(value)
	if n < minLen || n > maxLen {
		return fmt.Errorf("must contain from %d-%d characters", minLen, maxLen)
	}
	return nil
}

func ValidDateUsername(value string) error {
	if err := ValidDataString(value, 3, 100); err != nil {
		return err
	}

	if !isValidUsername(value) {
		return fmt.Errorf("username must contain letters, digits, underscore")
	}
	return nil
}

func ValidDateFullname(value string) error {
	if err := ValidDataString(value, 3, 100); err != nil {
		return err
	}

	if !isValidFullname(value) {
		return fmt.Errorf("Full name must contain letters")
	}
	return nil
}

func ValidDatePassword(value string) error {
	return ValidDataString(value, 6, 100)
}

func ValidDateEmail(value string) error {
	if err := ValidDataString(value, 3, 200); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("Is not a valid email address")
	}
	return nil
}
