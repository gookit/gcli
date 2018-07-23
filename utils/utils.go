package utils

import (
	"bytes"
	"encoding/json"
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

// ExecCommand
// cmdStr eg. "ls -al"
func ExecCommand(cmdStr string, shells ...string) (string, error) {
	shell := "/bin/sh"

	if len(shells) > 0 {
		shell = shells[0]
	}

	// 函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command(shell, "-c", cmdStr)

	// 读取io.Writer类型的cmd.Stdout，
	// 再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型
	var out bytes.Buffer
	cmd.Stdout = &out

	// Run执行c包含的命令，并阻塞直到完成。
	// 这里stdout被取出，cmd.Wait()无法正确获取 stdin,stdout,stderr，则阻塞在那了
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

// GetKeyMaxLen
// usage:
// utils.GetKeyMaxLen(map[string]string{"k1":"v1", "key2": "v2"}, 0)
func GetKeyMaxLen(kv map[string]interface{}, defLen int) (max int) {
	max = defLen

	for k, _ := range kv {
		kLen := len(k)

		if kLen > max {
			max = kLen
		}
	}

	return
}

// GetScreenSize
func GetScreenSize() (w int, h int) {
	return
}

// PrettyJson get pretty Json string
func PrettyJson(v interface{}) (string, error) {
	out, err := json.MarshalIndent(v, "", "    ")

	return string(out), err
}

// RenderTemplate
func RenderTemplate(input string, data interface{}, isFile ...bool) string {
	// use buffer receive rendered content
	var buf bytes.Buffer
	var isFilename bool

	if len(isFile) > 0 {
		isFilename = isFile[0]
	}

	t := template.New("cli")

	// don't escape content
	t.Funcs(template.FuncMap{"raw": func(s string) string {
		return s
	}})

	t.Funcs(template.FuncMap{"trim": func(s string) string {
		return strings.TrimSpace(string(s))
	}})

	// join strings
	t.Funcs(template.FuncMap{"join": func(ss []string, sep string) string {
		return strings.Join(ss, sep)
	}})

	// upper first char
	t.Funcs(template.FuncMap{"upFirst": func(s string) string {
		return UpperFirst(s)
	}})

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
