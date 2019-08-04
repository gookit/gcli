package gcli_test

import (
	"testing"

	"github.com/gookit/gcli/v2"
	"github.com/stretchr/testify/assert"
)

func TestInitStdApp(t *testing.T) {
	is := assert.New(t)

	gcli.InitStdApp(func(a *gcli.App) {
		a.Name = "test-name"
	})

	std := gcli.StdApp()

	is.Equal("test-name", std.Name)

	is.Empty(std.Commands())
}
