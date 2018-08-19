package utils

import (
	"strings"
)

// LcFirst char for given string
func LcFirst(s string) string {
	if len(s) == 0 {
		return s
	}

	f := s[0]
	isUpper := f >= 'A' && f <= 'Z'

	if isUpper {
		return strings.ToLower(string(f)) + string(s[1:])
	}

	return s
}

// UcFirst upper first char
func UcFirst(s string) string {
	if len(s) == 0 {
		return s
	}

	f := s[0]
	isLower := f >= 'a' && f <= 'z'

	if isLower {
		return strings.ToUpper(string(f)) + string(s[1:])
	}

	return s
}
