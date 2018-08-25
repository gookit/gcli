package emoji

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

var nameMatch = regexp.MustCompile(`(:\w+:)`)

// Emoji is alias of the GetByName()
func Emoji(name string) string {
	return GetByName(name)
}

// GetByName returns the unicode value for the given emoji name. If the
// specified emoji does not exist, will returns the input string.
func GetByName(name string) string {
	if val, ok := emojiMap[name]; ok {
		return val
	}

	return name
}

// Render a string, parse emoji name, returns rendered string.
// Usage:
// 	msg := Render("a :smile: message")
//	fmt.Println(msg)
func Render(str string) string {
	// not contains emoji name.
	if strings.IndexByte(str, ':') == -1 {
		return str
	}

	return nameMatch.ReplaceAllStringFunc(str, func(name string) string {
		return GetByName(str) // + " "
	})
}

// FromUnicode unicode string to emoji string
// Usage:
// 	emoji := FromUnicode("\U0001f496")
func FromUnicode(s string) string {
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

	ToUnicode('💖')

	return s
}

// ToUnicode unicode string to emoji string
// Usage:
// 	unicode := ToUnicode('💖')
//	fmt.Print(unicode) // "1f496"
//	// with prefix
// 	unicode := ToUnicode('💖', "\U000") // "\U0001f496"
//	fmt.Print(unicode) // "💖"
func ToUnicode(emoji rune, prefix ...string) string {
	code := strconv.FormatInt(int64(emoji), 16)

	if len(prefix) > 0 {
		return prefix[0] + code
	}

	return code
}

// Decode a string, convert unicode to emoji chat
// Usage:
// 	str := Decode("a msg [\u1f496]")
func Decode(s string) string {
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

	ToUnicode('💖')

	return s
}

// Encode a string, convert emoji chat to unicode string
func Encode(s string) string {
	rs := []rune(s)
	buf := new(bytes.Buffer)

	for _, r := range rs {
		if len(string(r)) == 4 { // is unicode emoji char
			code := strconv.FormatInt(int64(r), 16)
			buf.WriteString(`[\u` + code + `]`)
		} else {
			buf.WriteRune(r)
		}
	}

	return buf.String()
}
