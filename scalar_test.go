package perl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScalarString(t *testing.T) {
	i := NewInterpreter()
	s := argTypeMap(i, "한글")

	assert.Equal(t, "한글", s.String())
	assert.Equal(t, []byte("한글"), s.Bytes())
}

func TestScalarBytes(t *testing.T) {
	i := NewInterpreter()
	s := argTypeMap(i, []byte("한글"))

	assert.Equal(t, []byte("한글"), s.Bytes())
}

func TestScalarLatin1Bytes(t *testing.T) {
	i := NewInterpreter()
	s := argTypeMap(i, []byte("\xc9")) // é

	assert.Equal(t, []byte("\xc9"), s.Latin1Bytes())
}

func TestScalarUtf8Bytes(t *testing.T) {
	i := NewInterpreter()
	s := argTypeMap(i, []byte("\xc9")) // é

	assert.Equal(t, []byte("\xc3\x89"), s.Utf8Bytes())
}
