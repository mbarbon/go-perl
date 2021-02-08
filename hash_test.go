package perl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashBasic(t *testing.T) {
	i := NewInterpreter()
	r, err := i.EvalScalar(`{a => 3, b => 4}`)
	errPanic(err)
	h := r.Hash()

	assert.Equal(t, 3, h.FetchStringKey("a").Int())
	assert.Equal(t, 4, h.FetchStringKey("b").Int())
}
