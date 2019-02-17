package gcli_test

import (
	"testing"

	"github.com/gookit/gcli"
	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {
	is := assert.New(t)

	gcli.NewDefaultApp(func(a *gcli.App) {
		a.Name = "test-name"
	})

	is.Equal("test-name", gcli.DefaultApp.Name)
}
