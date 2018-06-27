package utils

import "fmt"

// The string flag list, implemented flag.Value interface
type StrList []string

func (s *StrList) String() string {
	return fmt.Sprintf("%s", *s)
}

func (s *StrList) Set(value string) error {
	*s = append(*s, value)
	return nil
}
