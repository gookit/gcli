package util

import (
	"io"
	"os"
)

func GetScreenSize() (w int, h int) {
	return
}

// 判断 w 是否为 stderr、stdout、stdin 三者之一
func IsConsole(out io.Writer) bool {
	o, ok := out.(*os.File)
	if !ok {
		return false
	}

	return o == os.Stdout || o == os.Stderr || o == os.Stdin
}
