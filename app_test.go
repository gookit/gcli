package cliapp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApp(t *testing.T) {
	is := assert.New(t)

	New(func(a *Application) {
		a.Name = "test-name"
	})

	is.Equal("test-name", App().Name)
}
