package httpapi

import (
	"regexp"
	"strings"
)

// emailPattern is a deliberately simple "valid enough for MVP" email check
// (docs/domain.md: guestEmail must pass simple MVP email validation).
var emailPattern = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

func isValidEmail(s string) bool {
	return emailPattern.MatchString(strings.TrimSpace(s))
}

func isBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}
