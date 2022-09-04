package interact

import "github.com/gookit/goutil/maputil"

// Collector information collector 信息收集者
type Collector struct {
	qs  []*Question
	ans maputil.Data
}

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) Run() error {
	return nil
}
