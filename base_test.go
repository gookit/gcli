package gcli_test

import (
	"testing"

	"github.com/gookit/gcli/v3"
)

func TestApp_On(t *testing.T) {
	cli := gcli.New()
	cli.ExitOnEnd = false

	args := []string{"top", "sub"}
	cli.Run(args)
}
