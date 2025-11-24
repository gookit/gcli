package title_test

import (
	"fmt"
	"testing"

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
	fmt.Println(content)
	assert.NotEmpty(t, content)

	// 测试顶部边框
	t2 := title.New("Top Border", title.WithBorderTop())
	content = t2.Render()
	if content == "" {
		t.Error("Expected non-empty content with top border")
	}

	// 测试底部边框
	t3 := title.New("Bottom Border", title.WithBorderBottom())
	content = t3.Render()
	if content == "" {
		t.Error("Expected non-empty content with bottom border")
	}

	// 测试上下边框
	t4 := title.New("Top and Bottom Border", title.WithBorderBoth())
	content = t4.Render()
	if content == "" {
		t.Error("Expected non-empty content with top and bottom borders")
	}

	// 测试无边框
	t5 := title.New("No Border", title.WithoutBorder())
	content = t5.Render()
	if content == "" {
		t.Error("Expected non-empty content without border")
	}
}
