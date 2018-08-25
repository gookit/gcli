# record

- UTF-8字符集的编码范围 `\u0000 - \uFFFF`

## emoji表情

emoji表情采用的是 Unicode编码，Emoji就是一种在Unicode位于 `\u1F601-\u1F64F`区段的字符。

显然超过了目前常用的UTF-8字符集的编码范围`\u0000-\uFFFF`。

```bash
#First, tan the following command to generate emojis.txt
wget -qO- https://unicode.org/Public/emoji/11.0/emoji-test.txt | cut -f 1 -d ' ' | sort -u | sed '/^[#0]/ d' | sed '/^\s*$/d' > /tmp/emojis.txt
```

## data links:

- https://www.unicode.org/Public/emoji/11.0/
- https://raw.githubusercontent.com/github/gemoji/master/db/emoji.json
- https://github.com/muan/emoji/blob/gh-pages/javascripts/emojilib/emojis.json

## sites

- http://emoji.muan.co/
- http://unicode-table.com/
- https://unicode-table.com/cn/sets/arrows-symbols/

## refer

- https://github.com/unicode-table/unicode-table-data 
- https://github.com/muan/emoji
- https://github.com/kyokomi/generateEmojiCodeMap
