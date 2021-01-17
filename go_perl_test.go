package perl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	i := NewInterpreter()

	listSeparator := i.Scalar(`"`)
	assert.Equal(t, " ", listSeparator.String())

	subscriptSeparator := i.Scalar(`;`)
	assert.Nil(t, subscriptSeparator)
	subscriptSeparator = i.CreateScalar(`;`)
	assert.Equal(t, "\x1C", subscriptSeparator.String())
}
