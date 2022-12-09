package interact_test

import (
	"testing"

	"github.com/gookit/gcli/v3/interact"
	"github.com/gookit/gcli/v3/interact/cparam"
	"github.com/gookit/goutil/testutil/assert"
)

func TestCollector_Run(t *testing.T) {
	c := interact.NewCollector()
	err := c.AddParams(
		cparam.NewStringParam("title", "title name"),
		cparam.NewChoiceParam("projects", "select projects"),
	)

	assert.NoErr(t, err)
}
