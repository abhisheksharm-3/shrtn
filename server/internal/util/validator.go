package util

import (
	"net/url"
	"regexp"
)

// ValidateURL checks if a string is a valid URL
func ValidateURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// ValidateShortCode checks if a string is a valid short code
func ValidateShortCode(code string) bool {
	// Short code should be alphanumeric and between 3-10 characters
	match, _ := regexp.MatchString("^[a-zA-Z0-9]{3,10}$", code)
	return match
}
