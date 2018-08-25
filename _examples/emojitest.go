package main

import (
	"fmt"
	"github.com/gookit/cliapp/show/emoji"
	"regexp"
	"strconv"
	"strings"
	"github.com/gookit/cliapp/show/symbols"
)

// go run ./_examples/test.go
func main() {
	fmt.Println(symbols.LEFT, emoji.BOX, "\xe2\x8c\x9a", emoji.HEART, "\u2764", "END")
	fmt.Println(emoji.HEART, "🚻", "\U0001f44d", "\U0001F17E", "\U00000038\U000020e3", "\U0001f4af")

	fmt.Println("\u2601\U000FE001", Emojitize("hello :snake: emoji :es:"), string(emoji.HEART1))

	fmt.Println(UnicodeEmojiCode(emoji.HEART), "\U0001F194", UnicodeEmojiDecode("\U0001f496"))

	ns := emoji.ToUnicode(emoji.HEART1)
	fmt.Println(ns, "👩🏾👩🏽", "\U0001F469\U0001F3FD", "\U0001f170")
}

var emojiMap = map[string]string{
	":strawberry:":  "\U0001f353",
	":worried:":     "\U0001f61f",
	":circus_tent:": "\U0001f3aa",
	":weary:":       "\U0001f629",
	":bathtub:":     "\U0001f6c1",
	":snake:":       "\U0001f40d",
	":grin:":        "\U0001f601",
	":symbols:":     "\U0001f523",
	":jp:":          "\U0001f1ef\U0001f1f5",
	":es:":          "\U0001f1ea\U0001f1f8",
}

// Emoji returns the unicode value for the given emoji. If the
// specified emoji does not exist, Emoji() returns the empty string.
func Emoji(emoji string) string {
	val, ok := emojiMap[emoji]
	if !ok {
		return emoji
	}
	return val
}

var reg = regexp.MustCompile("(:\\w+:)")

// Emojitize takes in a string with emojis specified in it, and returns
// a string with every emoji place holder replaced with it's unicode value
// (unless it could not be found, in which case it is let alone).
func Emojitize(emojis string) string {
	return reg.ReplaceAllStringFunc(emojis, func(str string) string {
		return Emoji(str) // + " "
	})
}

// 表情解码
func UnicodeEmojiDecode(s string) string {
	// emoji表情的数据表达式
	re := regexp.MustCompile("\\[[\\\\u0-9a-zA-Z]+\\]")
	// 提取emoji数据表达式
	reg := regexp.MustCompile("\\[\\\\u|]")
	src := re.FindAllString(s, -1)
	for i := 0; i < len(src); i++ {
		e := reg.ReplaceAllString(src[i], "")
		p, err := strconv.ParseInt(e, 16, 32)
		if err == nil {
			s = strings.Replace(s, src[i], string(rune(p)), -1)
		}
	}
	return s
}

// 表情转换
func UnicodeEmojiCode(s string) string {
	ret := ""
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if len(string(rs[i])) == 4 {
			u := `[\u` + strconv.FormatInt(int64(rs[i]), 16) + `]`
			ret += u

		} else {
			ret += string(rs[i])
		}
	}
	return ret
}
