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

func TestScalarType(t *testing.T) {
	testScalarType(t, `undef`, Undef)
	testScalarType(t, `$a = 1; undef $a; $a`, Undef)
	testScalarType(t, `$a = 1; $a`, Int)
	testScalarType(t, `$a = 0x8000000000000000; $a`, UInt)
	testScalarType(t, `$a = 1.3; $a`, Float)
	testScalarType(t, `$a = 'abc'; $a`, String)
	testScalarType(t, `$a = '한글'; $a`, String)
	testScalarType(t, `$a = \1; $a`, ScalarRef)
	testScalarType(t, `$a = []; $a`, ArrayRef)
	testScalarType(t, `$a = {}; $a`, HashRef)
	testScalarType(t, `$a = sub { $b }; $a`, CodeRef)
	testScalarType(t, `$a = qr/abc/; $a`, RegexpRef)
}

func TestScalarCode(t *testing.T) {
	i := NewInterpreter()
	s, err := i.EvalScalar(`sub { 'abc' }`)
	errPanic(err)

	sub := s.Code()
	assert.NotNil(t, sub)

	r, err := sub.CallScalar()
	errPanic(err)

	assert.Equal(t, "abc", r.String())
}

func testScalarType(t *testing.T, code string, expected ScalarType) {
	i := NewInterpreter()
	s, err := i.EvalScalar(code)
	errPanic(err)
	assert.Equal(t, expected, s.Type())
}
