package utils

import (
	"os/exec"
	"bytes"
	"encoding/json"
)

// ExecCommand
// cmdStr eg. "ls -al"
func ExecCommand(cmdStr string, shells ...string) (string, error) {
	shell := "/bin/bash"

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
