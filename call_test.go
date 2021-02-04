package perl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListEval(t *testing.T) {
	i := NewInterpreter()

	{
		res, err := i.EvalList(`"a", "b", "c"`)
		assert.NoError(t, err)

		assert.Equal(t, 3, len(res))
		assert.Equal(t, "a", res[0].String())
		assert.Equal(t, "b", res[1].String())
		assert.Equal(t, "c", res[2].String())
	}

	{
		res, err := i.EvalList(`()`)
		assert.NoError(t, err)

		assert.Equal(t, 0, len(res))
	}
}

func TestScalarEval(t *testing.T) {
	i := NewInterpreter()

	res, err := i.EvalScalar(`"a", "b", "c"`)
	assert.NoError(t, err)

	assert.Equal(t, "c", res.String())
}

func TestListScall(t *testing.T) {
	i := NewInterpreter()

	{
		err := i.EvalVoid(`sub test1 { "a", "b", "c" }`)
		assert.NoError(t, err)

		sub := i.Sub("test1")
		res, err := sub.CallList()
		assert.NoError(t, err)

		assert.Equal(t, 3, len(res))
		assert.Equal(t, "a", res[0].String())
		assert.Equal(t, "b", res[1].String())
		assert.Equal(t, "c", res[2].String())
	}

	{
		err := i.EvalVoid(`sub test2 { }`)
		assert.NoError(t, err)

		sub := i.Sub("test2")
		res, err := sub.CallList()
		assert.NoError(t, err)

		assert.Equal(t, 0, len(res))
	}
}

func TestScalarCall(t *testing.T) {
	i := NewInterpreter()

	err := i.EvalVoid(`sub test { "a", "b", "c" }`)
	assert.NoError(t, err)

	sub := i.Sub("test")
	res, err := sub.CallScalar()
	assert.NoError(t, err)

	assert.Equal(t, "c", res.String())
}
