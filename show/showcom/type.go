package showcom

// OverflowFlag for handling content overflow
type OverflowFlag uint8

// OverflowFlag values
const (
	OverflowAuto OverflowFlag = iota // auto: default is cut
	OverflowCut                      // 截断
	OverflowWrap                     // 换行
)
