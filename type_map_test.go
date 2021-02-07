package perl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeMapInts(t *testing.T) {
	testIntConversion(t, int(77), 77)
	testIntConversion(t, int16(77), 77)
	testIntConversion(t, int32(77), 77)
	testIntConversion(t, int64(77), 77)
}

func TestTypeMapUInts(t *testing.T) {
	testUIntConversion(t, uint(77), 77)
	testUIntConversion(t, uint16(77), 77)
	testUIntConversion(t, uint32(77), 77)
	testUIntConversion(t, uint64(77), 77)
}

func TestTypeMapFloat(t *testing.T) {
	testFloatConversion(t, float32(4.3), float64(float32(4.3)))
	testFloatConversion(t, float64(4.3), 4.3)
}

func TestTypeMapString(t *testing.T) {
	testStringConversion(t, "", 0, 0)
	testStringConversion(t, "abc", 0, 3)
	testStringConversion(t, "\xff\x41", 0, 2)
	testStringConversion(t, "한글", 1, 2)
}

func TestTypeMapBytes(t *testing.T) {
	testStringConversion(t, []byte(""), 0, 0)
	testStringConversion(t, []byte("abc"), 0, 3)
	testStringConversion(t, []byte("\xff\x41"), 0, 2)
	testStringConversion(t, []byte("한글"), 0, 6)
}

func TestTypeMapUtf8String(t *testing.T) {
	testStringConversion(t, Utf8String(""), 1, 0)
	testStringConversion(t, Utf8String("abc"), 1, 3)
	testStringConversion(t, Utf8String("\xff\x41"), 2, -1)
	testStringConversion(t, Utf8String("한글"), 1, 2)
}

func TestTypeMapUtf8Bytes(t *testing.T) {
	testStringConversion(t, Utf8Bytes([]byte("")), 1, 0)
	testStringConversion(t, Utf8Bytes([]byte("abc")), 1, 3)
	testStringConversion(t, Utf8Bytes([]byte("\xff\x41")), 2, -1)
	testStringConversion(t, Utf8Bytes([]byte("한글")), 1, 2)
}

func TestTypeMapByteString(t *testing.T) {
	testStringConversion(t, ByteString(""), 0, 0)
	testStringConversion(t, ByteString("abc"), 0, 3)
	testStringConversion(t, ByteString("\xff\x41"), 0, 2)
	testStringConversion(t, ByteString("한글"), 0, 6)
}

func testIntConversion(t *testing.T, goValue interface{}, expected int64) {
	i := NewInterpreter()
	s, err := toPerlArgScalar(i, goValue)
	errPanic(err)
	r := newScalarFromMortal(i, s)

	assert.Equal(t, expected, r.Int64())
}

func testUIntConversion(t *testing.T, goValue interface{}, expected uint64) {
	i := NewInterpreter()
	s, err := toPerlArgScalar(i, goValue)
	errPanic(err)
	r := newScalarFromMortal(i, s)

	assert.Equal(t, expected, r.UInt64())
}

func testFloatConversion(t *testing.T, goValue interface{}, expected float64) {
	i := NewInterpreter()
	s, err := toPerlArgScalar(i, goValue)
	errPanic(err)
	r := newScalarFromMortal(i, s)

	assert.Equal(t, expected, r.Float64())
}

func testStringConversion(t *testing.T, goValue interface{}, expectedUtf8, expectedLength int) {
	i := NewInterpreter()

	eval(i, `sub is_utf8 { length($_[0]), utf8::is_utf8($_[0]) ? "yes" : "no" }`)
	isUtf8 := i.Sub("is_utf8")

	r, err := isUtf8.CallList(goValue)

	switch expectedUtf8 {
	case 0:
		assert.Equal(t, "no", r[1].String())
		assert.Equal(t, expectedLength, r[0].Int())
	case 1:
		assert.Equal(t, "yes", r[1].String())
		assert.Equal(t, expectedLength, r[0].Int())
	case 2:
		assert.EqualError(t, err, ErrInvalidUtf8String.Error())
	default:
		panic("Not reached")
	}
}
