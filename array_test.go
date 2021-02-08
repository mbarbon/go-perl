package perl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayBasic(t *testing.T) {
	i := NewInterpreter()
	r, err := i.EvalScalar(`["a", 3]`)
	errPanic(err)
	a := r.Array()

	assert.Equal(t, "a", a.Fetch(0).String())
	assert.Equal(t, 3, a.Fetch(1).Int())
}
