package flags

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgsParser_Parse(t *testing.T) {
	at := assert.New(t)

	sample := "-n 10 --name tom --debug -aux --age 24 arg0 arg1"
	args := strings.Split(sample, " ")
	p := &ArgsParser{}
	p.Parse(args)

	at.Equal("[arg0 arg1]", fmt.Sprint(p.Args()))
	assertOptions := func(str string) {
		at.Contains(str, `"n":"10"`)
		at.Contains(str, `"name":"tom"`)
		at.Contains(str, `"debug":true`)
		at.Contains(str, `"age":"24"`)
		at.Contains(str, `"a":true`)
		at.Contains(str, `"u":true`)
		at.Contains(str, `"x":true`)
	}
	str := p.OptsString()
	assertOptions(str)

	// define bool options
	sample = "-n 10 --name tom --debug arg0 -aux --age 24 arg1"
	args = strings.Split(sample, " ")

	// 加上 []string{"debug"} 解析器就能正确分别 "--debug arg0"
	p = ParseArgs(args, []string{"debug"}, nil)
	at.Equal("[arg0 arg1]", fmt.Sprint(p.Args()))
	str = p.OptsString()
	assertOptions(str)

	// define array options
	sample = "-n 10 --name tom --name john --debug false -aux --age 24 arg0 arg1"
	args = strings.Split(sample, " ")
	// 加上 []string{"name"} 解析器就能正确解析 "--name --name john"
	p = ParseArgs(args, nil, []string{"name"})
	at.Equal("[arg0 arg1]", fmt.Sprint(p.Args()))
	str = p.OptsString()
	at.Contains(str, `"n":"10"`)
	at.Contains(str, `"name":[]string{"tom", "john"}`)
	at.Contains(str, `"debug":false`)
	at.Contains(str, `"age":"24"`)
	at.Contains(str, `"a":true`)
	at.Contains(str, `"u":true`)
	at.Contains(str, `"x":true`)

	// define bool and array options
	sample = "-n 10 --name tom --name john --debug arg0 -aux --age 24 arg1"
	args = strings.Split(sample, " ")

	p = ParseArgs(args, []string{"debug"}, []string{"name"})
	at.Equal("[arg0 arg1]", fmt.Sprint(p.Args()))
	str = p.OptsString()
	at.Contains(str, `"n":"10"`)
	at.Contains(str, `"name":[]string{"tom", "john"}`)
	at.Contains(str, `"debug":true`)
	at.Contains(str, `"age":"24"`)
	at.Contains(str, `"a":true`)
	at.Contains(str, `"u":true`)
	at.Contains(str, `"x":true`)
}
