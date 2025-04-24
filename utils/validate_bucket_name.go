package utils

import (
	"net"
	"regexp"
	"strings"
)

// Bucket naming rules: https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucketnamingrules.html
var bucketNameRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9\-\.]*[a-z0-9]$`)

// ValidateBucketName checks if the provided bucket name is valid.
func ValidateBucketName(name string) bool {
	return validateName(name, bucketNameRegex, 3, 63)
}

// validateName checks if the provided name matches the given regex and length constraints.
func validateName(name string, regex *regexp.Regexp, minLength, maxLength int) bool {
	if len(name) < minLength || len(name) > maxLength {
		return false
	}

	if !regex.MatchString(name) || isIPAddress(name) || hasConsecutivePeriodsOrDashes(name) {
		return false
	}

	if hasInvalidPrefix(name) || hasInvalidSuffix(name) {
		return false
	}

	return true
}

// isIPAddress checks if the given name is a valid IP address.
func isIPAddress(name string) bool {
	return net.ParseIP(name) != nil
}

// hasConsecutivePeriodsOrDashes checks for consecutive periods (..) or consecutive dashes (--) in the name.
func hasConsecutivePeriodsOrDashes(name string) bool {
	for i := 1; i < len(name); i++ {
		if (name[i] == '.' || name[i] == '-') && (name[i-1] == '.' || name[i-1] == '-') {
			return true
		}
	}
	return false
}

// hasInvalidPrefix checks for invalid prefixes.
func hasInvalidPrefix(name string) bool {
	invalidPrefixes := []string{"xn--", "sthree-", "sthree-configurator", "amzn-s3-demo-"}
	for _, prefix := range invalidPrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}

// hasInvalidSuffix checks for invalid suffixes.
func hasInvalidSuffix(name string) bool {
	invalidSuffixes := []string{"-s3alias", "--ol-s3", ".mrap", "--x-s3"}
	for _, suffix := range invalidSuffixes {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}
	return false
}
