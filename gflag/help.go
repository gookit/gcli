package gflag

import (
	"bytes"
	"flag"
	"fmt"
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/goutil/cflag"
	"github.com/gookit/goutil/strutil"
)

/***********************************************************************
 * Flags:
 * - render help message
 ***********************************************************************/

// SetHelpRender set the raw *flag.FlagSet.Usage
func (p *Parser) SetHelpRender(fn func()) {
	p.fSet.Usage = fn
}

// PrintHelpPanel for all options to the gf.out
func (p *Parser) PrintHelpPanel() {
	color.Fprint(p.out, p.String())
}

// String for all flag options
func (p *Parser) String() string {
	return p.BuildHelp()
}

// BuildHelp string for all flag options
func (p *Parser) BuildHelp() string {
	if p.buf == nil {
		p.buf = new(bytes.Buffer)
	}

	// repeat call the method
	if p.buf.Len() < 1 {
		p.buf.WriteString("Options:\n")
		p.buf.WriteString(p.BuildOptsHelp())
		p.buf.WriteByte('\n')

		if p.HasArgs() {
			p.buf.WriteString("Arguments:\n")
			p.buf.WriteString(p.BuildArgsHelp())
			p.buf.WriteByte('\n')
		}
	}

	return p.buf.String()
}

// BuildOptsHelp string.
func (p *Parser) BuildOptsHelp() string {
	var sb strings.Builder

	p.fSet.VisitAll(func(f *flag.Flag) {
		line := p.formatOneFlag(f)
		if line != "" {
			sb.WriteString(line)
			sb.WriteByte('\n')
		}
	})

	return sb.String()
}

func (p *Parser) formatOneFlag(f *flag.Flag) (s string) {
	// Skip render:
	// - opt is not exists(Has ensured that it is not a short name)
	// - it is hidden flag option
	// - flag desc is empty
	opt, has := p.opts[f.Name]
	if !has || opt.Hidden {
		return
	}

	var fullName string
	name := f.Name
	// eg: "-V, --version" length is: 13
	nameLen := p.names[name]
	// display description on new line
	descNl := p.cfg.DescNewline

	var nlIndent string
	if descNl {
		nlIndent = "\n        "
	} else {
		nlIndent = "\n      " + strings.Repeat(" ", p.optMaxLen)
	}

	// add prefix '-' to option
	fullName = cflag.AddPrefixes2(name, opt.Shorts, true)
	s = fmt.Sprintf("  <info>%s</>", fullName)

	// - build flag type info
	typeName, desc := flag.UnquoteUsage(f)
	// typeName: option value data type: int, string, ..., bool value will return ""
	if !p.cfg.WithoutType && len(typeName) > 0 {
		typeLen := len(typeName) + 1
		if !descNl && nameLen+typeLen > p.optMaxLen {
			descNl = true
		} else {
			nameLen += typeLen
		}

		s += fmt.Sprintf(" <magenta>%s</>", typeName)
	}

	if descNl {
		s += nlIndent
	} else {
		// padding space to optMaxLen width.
		if padLen := p.optMaxLen - nameLen; padLen > 0 {
			s += strings.Repeat(" ", padLen)
		}
		s += "    "
	}

	// --- build description
	if desc == "" {
		desc = defaultDesc
	} else {
		desc = strings.Replace(strutil.UpperFirst(desc), "\n", nlIndent, -1)
	}

	s += getRequiredMark(opt.Required) + desc

	// ---- append default value
	if isZero, isStr := cflag.IsZeroValue(f, f.DefValue); !isZero {
		if isStr {
			s += fmt.Sprintf(" (default <magentaB>%q</>)", f.DefValue)
		} else {
			s += fmt.Sprintf(" (default <magentaB>%v</>)", f.DefValue)
		}
	}

	// arrayed, repeatable
	if _, ok := f.Value.(cflag.RepeatableFlag); ok {
		s += " <cyan>(repeatable)</>"
	}
	return s
}
