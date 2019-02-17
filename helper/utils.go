package helper

import (
	"bytes"
	"os/exec"
	"strings"
	"text/template"

	"github.com/gookit/goutil/strutil"
)

// ExecCommand alias of the ShellExec
func ExecCommand(cmdStr string, dirAndShell ...string) (string, error) {
	return ShellExec(cmdStr, dirAndShell...)
}

// ShellExec exec a CLI command by shell and return output.
// Usage:
// 	ShellExec("ls -al")
// 	ShellExec("ls -al", "/usr/lib")
// 	ShellExec("ls -al", "/usr/lib", "/bin/zsh")
func ShellExec(cmdStr string, dirAndShell ...string) (string, error) {
	var workDir string
	shell := "/bin/sh"

	// if has more args
	if ln := len(dirAndShell); ln > 0 {
		workDir = dirAndShell[0]
		if ln > 1 {
			shell = dirAndShell[1]
		}
	}

	// create a new Cmd instance
	cmd := exec.Command(shell, "-c", cmdStr)
	if workDir != "" {
		cmd.Dir = workDir
	}

	bs, err := cmd.Output()
	return string(bs), err
}

// GetScreenSize for current console terminal
func GetScreenSize() (w int, h int) {
	return
}

// RenderText render text template with data
func RenderText(input string, data interface{}, fns template.FuncMap, isFile ...bool) string {
	// use buffer receive rendered content
	var buf bytes.Buffer

	t := template.New("cli")
	t.Funcs(template.FuncMap{
		// don't escape content
		"raw": func(s string) string {
			return s
		},
		"trim": strings.TrimSpace,
		// join strings. usage {{ join .Strings ","}}
		"join": func(ss []string, sep string) string {
			return strings.Join(ss, sep)
		},
		// lower first char
		"lcFirst": strutil.LowerFirst,
		// upper first char
		"ucFirst": strutil.UpperFirst,
	})

	// custom add template functions
	if len(fns) > 0 {
		t.Funcs(fns)
	}

	if len(isFile) > 0 && isFile[0] {
		template.Must(t.ParseFiles(input))
	} else {
		template.Must(t.Parse(input))
	}

	if err := t.Execute(&buf, data); err != nil {
		panic(err)
	}

	return buf.String()
}
