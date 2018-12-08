package gcli_test

import (
	"github.com/gookit/gcli"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApp(t *testing.T) {
	is := assert.New(t)

	gcli.NewApp(func(a *gcli.App) {
		a.Name = "test-name"
	})

	is.Equal("test-name", gcli.Instance().Name)
}
