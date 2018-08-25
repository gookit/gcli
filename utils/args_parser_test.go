package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestArgsParser_Parse(t *testing.T) {
	at := assert.New(t)
	p := ArgsParser{}
	at.NotNil(p)

	str := "-n 10 --name tom --debug --age 24 arg0 arg1"
	ret := p.Parse(strings.Split(str, " "))

	fmt.Printf("%#v\n", ret)
}
