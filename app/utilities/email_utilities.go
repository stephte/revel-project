package utilities

import (
	"github.com/revel/revel"
	"regexp"
)

func IsValidEmail(email string) bool {
	pattern := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\\.[a-zA-Z0-9]{2,61}(?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])*$")

	return pattern.MatchString(email)
}

func FilterValidEmails(emails []string) []string {
	rv := []string{}
	for _, email := range emails {
		if IsValidEmail(email) {
			rv = append(rv, email)
		} else {
			revel.AppLog.Warnf("Invalid email address: %s", email)
		}
	}

	return rv
}
