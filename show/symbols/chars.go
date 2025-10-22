package symbols

// links:
//
//	http://cn.piliapp.com/symbol/
//
// 卍 卐 ■ ▶ ☐☑☒ ❖
const (
	OK  = '✔'
	NO  = '✘'
	PEN = '✎'

	Center  rune = '●'
	Square  rune = '■'
	Square1 rune = '▇'
	Square2 rune = '▉'
	Square3 rune = '░'
	Square4 rune = '▒'
	Square5 rune = '▢'

	HEART  = '❤'
	HEART1 = '♥'
	SMILE  = '☺'

	FLOWER = '✿'
	MUSIC  = '♬'

	// UP ☚ ☜ ☛ ☞
	UP     = '⇧'
	DOWN   = '⇩'
	LEFT   = '⇦'
	RIGHT  = '⇨'
	SEARCH = ''

	// ❝❞❛❜
	// ⌜⌝⌞⌟
	// ▶➔➙➛➜➞➟➠➡➢➣➥➦➧➨➩➪➫➬➭➮➯➱➵

	MALE   = '♂'
	FEMALE = '♀'

	SUN   = '☀'
	STAR  = '★'
	SNOW  = '❈'
	CLOUD = '☁'

	ENTER = '⌥'

	Star   rune = '*'
	Plus   rune = '+'
	Well   rune = '#'
	Equal  rune = '='
	Equal1 rune = '═'
	Space  rune = ' '

	Underline  rune = '_'
	LeftArrow  rune = '<'
	RightArrow rune = '>'
)

// style1:
// ╭─────╮
// │  hi │
// ╰─────╯
// style2:
//  ┌────┐
//  │ hi │
//  └────┘
const (
	// Hyphen Minus
	Hyphen   rune = '-' // eg: -------
	CNHyphen rune = '—' // eg: —————
	Hyphen2  rune = '─' // eg: ────

	VLine     rune = '|'
	VLineFull rune = '│'

	LeftTop1 = '╭'
	LeftTop2 = '┌'

	RightTop1 = '╮'
	RightTop2 = '┐'

	LeftBottom1 = '╰'
	LeftBottom2 = '└'

	RightBottom1 = '╯'
	RightBottom2 = '┘'

	// TChar eg TChar + Hyphen2: ──┬──
	TChar rune = '┬'
	// CCChar criss-cross
	CCChar rune = '┼'
)
