package utils

import (
	"bytes"
	"github.com/gookit/goutil/strUtil"
	"os/exec"
	"strings"
	"text/template"
)

// ExecCommand alias of the ShellExec
func ExecCommand(cmdStr string, dirAndShell ...string) (string, error) {
	return ShellExec(cmdStr, dirAndShell...)
}

// ShellExec exec a CLI command by shell and return output.
// usage:
// 	utils.ShellExec("ls -al")
// 	utils.ShellExec("ls -al", "/usr/lib")
// 	utils.ShellExec("ls -al", "/usr/lib", "/bin/zsh")
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

// RenderTemplate render text template with data
func RenderTemplate(input string, data interface{}, isFile ...bool) string {
	// use buffer receive rendered content
	var buf bytes.Buffer
	var isFilename bool

	if len(isFile) > 0 {
		isFilename = isFile[0]
	}

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
		"lcFirst": strUtil.LowerFirst,
		// upper first char
		"ucFirst": strUtil.UpperFirst,
	})

	if isFilename {
		template.Must(t.ParseFiles(input))
	} else {
		template.Must(t.Parse(input))
	}

	if err := t.Execute(&buf, data); err != nil {
		panic(err)
	}

	return buf.String()
}
