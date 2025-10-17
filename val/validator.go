package val

import (
	"fmt"
	"net/mail"
	"regexp"

	"github.com/JeongWoo-Seo/simpleBank/util"
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

func ValidDataID(value int64) error {
	if value <= 0 {
		return fmt.Errorf("must be a positive integer")
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
		return fmt.Errorf("is not a valid email address")
	}
	return nil
}

func ValidDateCurrency(value string) error {
	if !util.IsSupportedCurrency(value) {
		return fmt.Errorf("unsupport currency")
	}
	return nil
}

func ValidDataEmailID(value int64) error {
	if value <= 0 {
		return fmt.Errorf("must be a positive integer")
	}
	return nil
}

func ValidDataSecretCode(value string) error {
	return ValidDataString(value, 32, 128)
}

func ValidDatePageID(value int32) error {
	if value <= 0 {
		return fmt.Errorf("must be a positive integer")
	}
	return nil
}

func ValidDatePageSize(value int32) error {
	if value < 5 || value > 10 {
		return fmt.Errorf("page sizd is 5 ~ 10")
	}
	return nil
}
