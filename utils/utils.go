package utils

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

// Go is a basic promise implementation: it wraps calls a function in a goroutine
// and returns a channel which will later return the function's return value.
// from beego/bee
func Go(f func() error) chan error {
	ch := make(chan error)
	go func() {
		ch <- f()
	}()
	return ch
}

// ExecCmd exec a CLI command and return output.
func ExecCmd(cmdStr string, dirAndShell ...string) (string, error) {
	return ExecCommand(cmdStr, dirAndShell...)
}

// ExecCommand exec a CLI command and return output.
// usage:
// 	utils.ExecCommand("ls -al")
// 	utils.ExecCommand("ls -al", "/usr/lib")
// 	utils.ExecCommand("ls -al", "/usr/lib", "/bin/zsh")
func ExecCommand(cmdStr string, dirAndShell ...string) (string, error) {
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

	// rewrite cmd.Stdout to buffer
	out := new(bytes.Buffer)
	cmd.Stdout = out

	// Run执行命令，并阻塞直到完成
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return out.String(), nil
}

// GetCurShell get current used shell env file. eg "/bin/zsh" "/bin/bash"
func GetCurShell(onlyName bool) string {
	path, err := ExecCommand("echo $SHELL")

	if err != nil {
		return ""
	}

	path = strings.TrimSpace(path)
	if onlyName && len(path) > 0 {
		path = filepath.Base(path)
	}

	return path
}

// GetKeyMaxLen get key max length of the map
// usage:
// 	utils.GetKeyMaxLen(map[string]string{"k1":"v1", "key2": "v2"}, 0)
func GetKeyMaxLen(kv map[string]interface{}, defLen int) (max int) {
	max = defLen
	for k := range kv {
		kl := len(k)
		if kl > max {
			max = kl
		}
	}

	return
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
		"lcFirst": LcFirst,
		// upper first char
		"ucFirst": UcFirst,
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
