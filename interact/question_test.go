package interact_test

import (
	"testing"

	"github.com/gookit/gcli/v3/interact"
)

func TestQuestion_Run(t *testing.T) {
	q := interact.NewQuestion("your name")
	q.Run()
}
