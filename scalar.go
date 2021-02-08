package perl

/*
#include "go_perl.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// ScalarType represents the type of value held by a Scalar
type ScalarType int

const (
	// Undef is an undefined scalars
	Undef ScalarType = iota + 1
	// Int is a signed integer value
	Int
	// UInt is an unsigned integer value
	UInt
	// Float is a floating point value
	Float
	// String is a string value
	String
	// ScalarRef is a reference to a scalar value
	ScalarRef
	// ArrayRef is a reference to an array value
	ArrayRef
	// HashRef is a reference to an hash value
	HashRef
	// CodeRef is a reference to a code value
	CodeRef
	// RegexpRef is a reference to a regexp
	RegexpRef
	// Unknown is a reference to a scalar value
	Unknown
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

// Scalar returns the Scalar value wrapped by a scalar reference, or nil
func (s *Scalar) Scalar() *Scalar {
	if s.Type() != ScalarRef {
		return nil
	}
	return newScalarFromMortal(s.gpi, C.go_perl_svrv(s.gpi.pi, s.sv))
}

// Array returns the Array value wrapped by a scalar reference, or nil
func (s *Scalar) Array() *Array {
	if s.Type() != ArrayRef {
		return nil
	}
	return newArrayFromMortal(s.gpi, C.sv_to_av(C.go_perl_svrv(s.gpi.pi, s.sv)))
}

// Hash returns the Hash value wrapped by a scalar reference, or nil
func (s *Scalar) Hash() *Hash {
	if s.Type() != HashRef {
		return nil
	}
	return newHashFromMortal(s.gpi, C.sv_to_hv(C.go_perl_svrv(s.gpi.pi, s.sv)))
}

// Code returns the Sub value wrapped by a scalar reference, or nil
func (s *Scalar) Code() *Sub {
	if s.Type() != CodeRef {
		return nil
	}
	return newSubFromMortal(s.gpi, C.sv_to_cv(C.go_perl_svrv(s.gpi.pi, s.sv)))
}

// Type returns the type of value stored in this scalar
func (s *Scalar) Type() ScalarType {
	svType := C.go_perl_sv_type(s.gpi.pi, s.sv)
	svFlags := C.go_perl_decoded_sv_flags(s.gpi.pi, s.sv)

	if svFlags&C.GO_PERL_SVf_RV != 0 {
		rv := C.go_perl_svrv(s.gpi.pi, s.sv)
		rvType := C.go_perl_sv_type(s.gpi.pi, rv)
		rvFlags := C.go_perl_decoded_sv_flags(s.gpi.pi, rv)
		return refType(rvType, rvFlags)
	}

	switch svType {
	case C.SVt_NULL:
		return Undef
	case C.SVt_IV, C.SVt_NV, C.SVt_PV, C.SVt_PVIV, C.SVt_PVNV, C.SVt_PVMG:
		if svFlags&C.GO_PERL_SVf_NV != 0 {
			return Float
		} else if svFlags&C.GO_PERL_SVf_UV != 0 {
			return UInt
		} else if svFlags&C.GO_PERL_SVf_IV != 0 {
			return Int
		} else if svFlags&C.GO_PERL_SVf_PV != 0 {
			return String
		}
		return Undef
	default:
		return Unknown
	}
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

func refType(svType, svFlags C.int) ScalarType {
	switch svType {
	case C.SVt_NULL, C.SVt_IV, C.SVt_NV, C.SVt_PV, C.SVt_PVIV, C.SVt_PVNV, C.SVt_PVMG:
		return ScalarRef
	case C.SVt_REGEXP:
		return RegexpRef
	case C.SVt_PVAV:
		return ArrayRef
	case C.SVt_PVHV:
		return HashRef
	case C.SVt_PVCV:
		return CodeRef
	default:
		return Unknown
	}
}
