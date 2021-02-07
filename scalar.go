package perl

/*
#include "go_perl.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// Scalar wraps a Perl scalar value
type Scalar struct {
	gpi *Interpreter
	sv  *C.go_perl_sv
}

type scalarError struct {
	scalar *Scalar
}

var _ error = (*scalarError)(nil)

func newScalar(gpi *Interpreter, sv *C.go_perl_sv) *Scalar {
	return &Scalar{
		gpi: gpi,
		sv:  sv,
	}
}

func newScalarFromMortal(gpi *Interpreter, sv *C.go_perl_sv) *Scalar {
	s := newScalar(gpi, sv)
	C.go_perl_sv_refcnt_inc(gpi.pi, sv)
	s.setFinalizer()
	return s
}

// String returns the scalar value coerced to a string
func (s *Scalar) String() string {
	var len C.go_perl_strlen
	str := C.go_perl_svpv_utf8(s.gpi.pi, s.sv, &len)
	return C.GoStringN(str, C.int(len))
}

// Bytes returns the scalar value coerced to []byte
func (s *Scalar) Bytes() []byte {
	var len C.go_perl_strlen
	str := C.go_perl_svpv(s.gpi.pi, s.sv, &len)
	return C.GoBytes(unsafe.Pointer(str), C.int(len))
}

// Latin1Bytes encodes the scalar value as ISO-8859-1, if necessary
func (s *Scalar) Latin1Bytes() []byte {
	var len C.go_perl_strlen
	str := C.go_perl_svpv_bytes(s.gpi.pi, s.sv, &len)
	return C.GoBytes(unsafe.Pointer(str), C.int(len))
}

// Utf8Bytes encodes the scalar value as UTF-8, if necessary
func (s *Scalar) Utf8Bytes() []byte {
	var len C.go_perl_strlen
	str := C.go_perl_svpv_utf8(s.gpi.pi, s.sv, &len)
	return C.GoBytes(unsafe.Pointer(str), C.int(len))
}

// Int returns the current scalar value coerced to an int
func (s *Scalar) Int() int {
	return int(s.Int64())
}

// Int16 returns the current scalar value coerced to an int16
func (s *Scalar) Int16() int16 {
	return int16(s.Int64())
}

// Int32 returns the current scalar value coerced to an int32
func (s *Scalar) Int32() int32 {
	return int32(s.Int64())
}

// Int64 returns the current scalar value coerced to an int64
func (s *Scalar) Int64() int64 {
	return int64(C.go_perl_sviv(s.gpi.pi, s.sv))
}

// UInt returns the current scalar value coerced to an uint
func (s *Scalar) UInt() uint {
	return uint(s.UInt64())
}

// UInt16 returns the current scalar value coerced to an uint16
func (s *Scalar) UInt16() uint16 {
	return uint16(s.UInt64())
}

// UInt32 returns the current scalar value coerced to an uint32
func (s *Scalar) UInt32() uint32 {
	return uint32(s.UInt64())
}

// UInt64 returns the current scalar value coerced to an uint64
func (s *Scalar) UInt64() uint64 {
	return uint64(C.go_perl_svuv(s.gpi.pi, s.sv))
}

// Float32 returns the current scalar value coerced to a float32
func (s *Scalar) Float32() float32 {
	return float32(s.Float64())
}

// Float64 returns the current scalar value coerced to a float64
func (s *Scalar) Float64() float64 {
	return float64(C.go_perl_svnv(s.gpi.pi, s.sv))
}

func (s *Scalar) setFinalizer() {
	runtime.SetFinalizer(s, func(s *Scalar) {
		C.go_perl_sv_refcnt_dec(s.gpi.pi, s.sv)
	})
}

func (s *Scalar) asError() *scalarError {
	return &scalarError{scalar: s}
}

func (se *scalarError) Error() string {
	return se.scalar.String()
}
