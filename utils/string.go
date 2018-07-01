package utils

import (
	"strings"
)

// FindSimilar
func FindSimilar(input string, samples []string) []string {
	var ss []string
	// ins := strings.Split(input, "")

	// fmt.Print(input, ins)

	for _, str := range samples{
		if strings.Contains(str, input) {
			ss = append(ss, str)
		} else {
			// sns := strings.Split(str, "")
		}

		// max find four items
		if len(ss) == 4 {
			break
		}
	}

	// fmt.Println("found ", ss)

	return ss
}

// LowerFirst
func LowerFirst(s string) string {
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

// UpperFirst upper first char
func UpperFirst(s string) string {
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
