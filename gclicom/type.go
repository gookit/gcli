package gclicom

type TextPos uint8

const (
	TextPosLeft TextPos = iota
	TextPosCenter
	TextPosRight
)

const TextPosMiddle = TextPosCenter

type BorderPos uint8

const (
	BorderPosTop BorderPos = iota
	BorderPosBottom
	BorderPosLeft
	BorderPosRight
	BorderPosTB // Top & Bottom
	BorderPosLR // Left & Right
	BorderPosAll
)
