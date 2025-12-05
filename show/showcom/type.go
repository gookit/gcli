package showcom

// OverflowFlag for handling content overflow. 0=auto, 1=cut, 2=wrap
type OverflowFlag uint8

// OverflowFlag values
const (
	OverflowAuto OverflowFlag = iota // auto
	OverflowCut                      // 截断
	OverflowWrap                     // 换行
)
