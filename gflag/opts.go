package gflag

import (
	"flag"
	"fmt"
	"strings"

	"github.com/gookit/gcli/v3/helper"
	"github.com/gookit/goutil/cflag"
	"github.com/gookit/goutil/stdutil"
	"github.com/gookit/goutil/structs"
	"github.com/gookit/goutil/strutil"
)

/***********************************************************************
 * flag options metadata
 ***********************************************************************/

// FlagMeta alias of flag Option
type FlagMeta = Option

// Option define for a flag option
type Option struct {
	// go flag value
	flag *flag.Flag
	// Name of flag and description
	Name, Desc string
	// default value for the flag option
	DefVal any
	// wrapped the default value
	defVal *structs.Value
	// short names. eg: ["o", "a"]
	Shorts []string
	// EnvVar allow set flag value from ENV var
	EnvVar string

	// --- advanced settings

	// Hidden the option on help
	Hidden bool
	// Required the option is required
	Required bool
	// Validator support validate the option flag value
	Validator func(val string) error
	// TODO interactive question for collect value
	Question string
}

// newFlagOpt quick create an FlagMeta
func newFlagOpt(name, desc string, defVal any, shortcut string) *Option {
	return &Option{
		Name: name,
		Desc: desc,
		// other info
		DefVal: defVal,
		Shorts: strings.Split(shortcut, ","),
	}
}

func (m *Option) initCheck() string {
	if m.Desc != "" {
		desc := strings.Trim(m.Desc, "; ")
		if strings.ContainsRune(desc, ';') {
			// format: desc;required
			// format: desc;required;env TODO parse ENV var
			parts := strutil.SplitNTrimmed(desc, ";", 2)
			if ln := len(parts); ln > 1 {
				bl, err := strutil.Bool(parts[1])
				if err == nil && bl {
					desc = parts[0]
					m.Required = true
				}
			}
		}

		m.Desc = desc
	}

	// filter shorts
	if len(m.Shorts) > 0 {
		m.Shorts = cflag.FilterNames(m.Shorts)
	}
	return m.goodName()
}

// good name of the flag
func (m *Option) goodName() string {
	name := strings.Trim(m.Name, "- ")
	if name == "" {
		helper.Panicf("option flag name cannot be empty")
	}

	if !helper.IsGoodName(name) {
		helper.Panicf("option flag name '%s' is invalid, must match: %s", name, helper.RegGoodName)
	}

	// update self name
	m.Name = name
	return name
}

// Shorts2String join shorts to a string
func (m *Option) Shorts2String(sep ...string) string {
	if len(m.Shorts) == 0 {
		return ""
	}
	return strings.Join(m.Shorts, sepStr(sep))
}

// HelpName for show help
func (m *Option) HelpName() string {
	return cflag.AddPrefixes(m.Name, m.Shorts)
}

func (m *Option) helpNameLen() int {
	return len(m.HelpName())
}

// Validate the binding value
func (m *Option) Validate(val string) error {
	if m.Required && val == "" {
		return fmt.Errorf("flag '%s' is required", m.Name)
	}

	// call user custom validator
	if m.Validator != nil {
		return m.Validator(val)
	}
	return nil
}

// Flag value
func (m *Option) Flag() *flag.Flag {
	return m.flag
}

// DValue wrap the default value
func (m *Option) DValue() *stdutil.Value {
	if m.defVal == nil {
		m.defVal = &stdutil.Value{V: m.DefVal}
	}
	return m.defVal
}
