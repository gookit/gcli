package main

import (
	"fmt"

	"github.com/gookit/cliui/show/emoji"
	"github.com/gookit/cliui/show/symbols"
)

// go run ./_examples/test.go
func main() {
	fmt.Println(symbols.LEFT, emoji.BOX, "\xe2\x8c\x9a", emoji.HEART, "\u2764", "END")
	fmt.Println(emoji.HEART, "🚻", "\U0001f44d", "\U0001F17E", "\U00000038\U000020e3", "\U0001f4af")

	fmt.Println("\u2601\U000FE001", emoji.Render("hello :snake: emoji :car:"))

	fmt.Println(emoji.ToUnicode(emoji.HEART), "\U0001F194", emoji.Decode("\U0001f496"))

	ns := emoji.ToUnicode(emoji.HEART)
	fmt.Println(ns, "👩🏾👩🏽", "\U0001F469\U0001F3FD", "\U0001f170")

}
