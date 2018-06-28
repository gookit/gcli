package utils

import (
	"io"
	"os"
	"os/exec"
	"bytes"
	"runtime"
)

func GetScreenSize() (w int, h int) {
	return
}

// isMSys
func isMSys() bool {
	if len(os.Getenv("MSYSTEM")) > 0 { // msys 环境
		return true
	}

	return false
}

// IsWin
// linux windows darwin
func IsWin() bool {
	return runtime.GOOS == "windows"
}

// IsMac
func IsMac() bool {
	return runtime.GOOS == "darwin"
}

// IsLinux
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// 判断 w 是否为 stderr、stdout、stdin 三者之一
func IsConsole(out io.Writer) bool {
	o, ok := out.(*os.File)
	if !ok {
		return false
	}

	return o == os.Stdout || o == os.Stderr || o == os.Stdin
}

// ExecOsCommand
// cmdStr eg. "ls -al"
func ExecOsCommand(cmdStr string, shells ...string) (string, error) {
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

// FindSimilar
func FindSimilar(input string, samples []string) {

}

// GetKeyMaxLen
// usage: utils.GetKeyMaxLen(map[string]string{"k1":"v1", "key2": "v2"}, 0)
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
