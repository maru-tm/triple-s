package utils

import (
	"regexp"
	"strings"
)

// Naming Amazon S3 objects: https://docs.aws.amazon.com/AmazonS3/latest/userguide/object-keys.html
var objectKeyRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9\-_\.]*[a-z0-9]$`)

// ValidateObjectKey checks if the provided object key is valid.
func ValidateObjectKey(key string) bool {
	return validateKey(key, objectKeyRegex, 1, 1024)
}

func validateKey(key string, regex *regexp.Regexp, minLength, maxLength int) bool {
	if len(key) < minLength || len(key) > maxLength {
		return false
	}

	if !regex.MatchString(key) || requiresSpecialHandling(key) || containsInvalidCharacters(key) {
		return false
	}

	if strings.Contains(key, "..") || strings.EqualFold(key, "soap") || isIPAddress(key) {
		return false
	}

	if hasInvalidPrefix(key) || hasInvalidSuffix(key) {
		return false
	}

	return true
}

// requiresSpecialHandling checks if the key contains special characters.
func requiresSpecialHandling(key string) bool {
	specialChars := []string{"&", "$", "@", "=", ";", "/", ":", "+", " ", ",", "?"}
	for _, char := range specialChars {
		if strings.Contains(key, char) {
			return true
		}
	}
	return false
}

// containsInvalidCharacters checks for characters that should be avoided.
func containsInvalidCharacters(key string) bool {
	invalidChars := []string{"\\", "{", "}", "^", "%", "`", "]", "\"", ">", "[", "<", "#", "|", "~"}
	for _, char := range invalidChars {
		if strings.Contains(key, char) {
			return true
		}
	}
	return false
}
