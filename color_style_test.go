package cliapp

import (
	"testing"
	"regexp"
	"fmt"
)

func TestMatchTag(t *testing.T) {
	s := "<err>text</err>"

	reg := regexp.MustCompile(ColorTag)
	r := reg.FindAllString(s, -1)

	fmt.Printf("%+v\n", r)
}
