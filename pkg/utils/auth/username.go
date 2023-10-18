package auth

import "regexp"

func IsUsernameValid(u string) bool {
	pattern := "^[A-Za-z][A-Za-z_]*[A-Za-z]$"
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(u)
}
