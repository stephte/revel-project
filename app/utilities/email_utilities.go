package utilities

import (
	"regexp"
	"errors"
	"strings"
)

func HandleEmail(email string) (string, error) {
	pattern := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\\.[a-zA-Z0-9]{2,61}(?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])*$")

	if !pattern.MatchString(email) {
		return "", errors.New("Invalid Email Address")
	}

	return strings.ToLower(email), nil
}