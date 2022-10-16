package utilities

import (
	"golang.org/x/crypto/bcrypt"
	"fmt"
)

func CreateHash(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return hash, err
	}

	return hash, nil
}

func ComparePasswords(pwHash []byte, pwGiven string) bool {
	err := bcrypt.CompareHashAndPassword(pwHash, []byte(pwGiven))

	if err == nil {
		return true
	} else {
		// log.Error(err)
		fmt.Println(err)
		return false
	}
}

// TODO: add more rigorous password validation
func ValidatePassword(password string) bool {
	valid := false

	valid = (len([]rune(password)) >= 8) || valid

	return valid
}
