package cliapp_test

import (
	"github.com/gookit/cliapp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApp(t *testing.T) {
	is := assert.New(t)

	cliapp.New(func(a *cliapp.App) {
		a.Name = "test-name"
	})

	is.Equal("test-name", cliapp.Instance().Name)
}
