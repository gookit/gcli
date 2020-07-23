package gcli_test

import (
	"testing"

	"github.com/gookit/gcli/v2"
	"github.com/stretchr/testify/assert"
)

func TestGFlags_StrOpt(t *testing.T) {
	gf := gcli.NewGFlags("test")

	var str string
	gf.StrOpt(&str, &gcli.Meta{
		Name:   "test",
		UseFor: "test desc",
	})
}

func TestGFlags_FromStruct(t *testing.T) {
	gf := gcli.NewGFlags("test")

	type userOpts struct {
		Opt1 string `gcli:"name=opt;short=o,h;desc=message"`
	}

	err := gf.FromStruct(&userOpts{})

	assert.NoError(t, err)
}

func TestGFlags_PrintHelpPanel(t *testing.T) {
	gf := gcli.NewGFlags("test")

	testOpts := struct {
		opt1 int
		opt2 bool
		opt3 string
	}{}

	gf.StrOpt(&testOpts.opt3, &gcli.Meta{
		Name:   "test",
		UseFor: "test desc",
	})
}