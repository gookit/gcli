package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// RUN: go run ./_examples/testterm/term_size.go
func main() {
	fd := int(os.Stdout.Fd())

	// 检查是否为终端
	if !term.IsTerminal(fd) {
		fmt.Println("Not a terminal")
		return
	}

	width, height, err := term.GetSize(fd)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Terminal size: %d x %d\n", width, height)
}
