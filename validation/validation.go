package validation

import (
	"regexp"
	"time"
)

var emailRegexp *regexp.Regexp = regexp.MustCompile("^[^\\s-]+@[^\\s-]+\\.[^\\s-]+$")
var dateDDMMYYYYLayoutString string = "2/1/2006"
var hasSmallCharRegexp *regexp.Regexp = regexp.MustCompile("[a-z]+")
var hasCapitalCharRegexp *regexp.Regexp = regexp.MustCompile("[A-Z]+")
var hasNumberRegexp *regexp.Regexp = regexp.MustCompile("[0-9]+")

var PasswordRequirementString = "Requires at least 1 small case character, one capital case character, and one number"

func IsEmail(toTest string) bool {
	return emailRegexp.MatchString(toTest)
}

func IsDateDDMMYYYY(dateString string) bool {
	_, err := time.Parse(dateDDMMYYYYLayoutString, dateString)
	if err != nil {
		return false
	}
	return true
}

func hasEightCharsASmallACapitalAndOneNumberAtLeast(candidate string) bool {
	if len(candidate) < 8 {
		return false
	}
	failed := false
	var regexpsToUse []*regexp.Regexp = []*regexp.Regexp{hasSmallCharRegexp, hasCapitalCharRegexp, hasNumberRegexp}

	for _, usingRegexp := range regexpsToUse {
		ok := usingRegexp.MatchString(candidate)
		if !ok {
			failed = true
			break
		}
	}

	return !failed
}

func IsValidPassword(candidate string) bool {
	return hasEightCharsASmallACapitalAndOneNumberAtLeast(candidate)
}
