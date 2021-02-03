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

func TestBasicCall(t *testing.T) {
	i := NewInterpreter()

	err := i.EvalVoid(`sub test { $val = $_[0] }`)
	assert.NoError(t, err)

	testSub := i.Sub("test")

	err = testSub.CallVoid("hello, world")
	assert.NoError(t, err)

	val := i.Scalar("val")
	assert.Equal(t, "hello, world", val.String())
}

func TestFailedCall(t *testing.T) {
	i := NewInterpreter()

	err := i.EvalVoid(`sub test { die "" }`)
	assert.NoError(t, err)

	testSub := i.Sub("test")

	err = testSub.CallVoid()
	assert.EqualError(t, err, "Died at (eval 1) line 1.\n")
}
