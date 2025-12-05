package title_test

import (
	"fmt"
	"testing"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/show/title"
	"github.com/gookit/goutil/testutil/assert"
)

// ExampleBorderTop 演示顶部边框
func ExampleWithBorderTop() {
	title.New("Top Border Title", title.WithBorderTop()).Println()

	// 演示底部边框（默认）
	title.New("Bottom Border Title", title.WithBorderBottom()).Println()

	// 演示上下边框
	title.New("Top and Bottom Border Title", title.WithBorderBoth()).Println()

	// ExampleNoBorder 演示无边框
	title.New("No Border Title", title.WithoutBorder()).Println()
}

func TestTitleRender(t *testing.T) {
	// 测试基本标题渲染
	t1 := title.New("Test Title")
	content := t1.Render()
	assert.NotEmpty(t, content)
	color.Println(content)
	fmt.Println()

	// 测试顶部边框
	t2 := title.New("Top Border", title.WithBorderTop())
	content = t2.Render()
	assert.NotEmpty(t, content)
	color.Println(content)
	fmt.Println()

	// 测试底部边框
	t3 := title.New("Bottom Border Align Center", title.WithBorderBottom(), title.WithAlignCenter())
	content = t3.Render()
	assert.NotEmpty(t, content)
	color.Println(content)
	fmt.Println()

	// 测试上下边框
	t4 := title.New("Top and Bottom Border", title.WithBorderBoth())
	content = t4.Render()
	assert.NotEmpty(t, content)
	color.Println(content)
	fmt.Println()

	// 测试无边框
	t5 := title.New("No Border Align Right", title.WithoutBorder(), title.WithAlignRight())
	content = t5.Render()
	assert.NotEmpty(t, content)
	color.Println(content)
}
