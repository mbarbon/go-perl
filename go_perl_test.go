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

func TestSuccessfulEval(t *testing.T) {
	i := NewInterpreter()

	err := i.EvalVoid(`$foo = "test"`)
	assert.NoError(t, err)

	foo := i.Scalar("foo")
	assert.Equal(t, "test", foo.String())
}

func TestFailedEval(t *testing.T) {
	i := NewInterpreter()

	err := i.EvalVoid(`die ""`)
	assert.EqualError(t, err, "Died at (eval 1) line 1.\n")
}
