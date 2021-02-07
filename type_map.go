package perl

/*
#include "go_perl.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

// ErrInvalidUtf8String is returned when the conversion explictly asks for UTF-8, but the passed-in data is not valid UTF-8
var ErrInvalidUtf8String = errors.New("Invalid UTF-8 encoding when trying to convert a string")

// StringCoercion is an internal marker type used to indicate whether Go->Perl conversion should produce characters or bytes
type StringCoercion interface{}

type stringCoercion struct {
	bytes    []byte
	str      string
	wantUtf8 C.int
}

// Utf8Bytes can be used to mark a []byte to be converted to an UTF-8 Perl scalar
func Utf8Bytes(bytes []byte) StringCoercion {
	return stringCoercion{
		bytes:    bytes,
		wantUtf8: 2,
	}
}

// Utf8String can be used to mark a string to be converted to an UTF-8 Perl scalar
func Utf8String(str string) StringCoercion {
	return stringCoercion{
		str:      str,
		wantUtf8: 2,
	}
}

// ByteString can be used to mark a string to be converted to a non-UTF-8 Perl scalar
func ByteString(str string) StringCoercion {
	return stringCoercion{
		str:      str,
		wantUtf8: 0,
	}
}

func toPerlMortalString(gpi *Interpreter, val string, wantUtf8 C.int) *C.go_perl_sv {
	cstr := C.CString(val)
	defer C.free(unsafe.Pointer(cstr))
	return C.go_perl_new_mortal_sv_string(gpi.pi, cstr, C.go_perl_strlen(len(val)), wantUtf8)
}

func toPerlMortalBinary(gpi *Interpreter, val []byte, wantUtf8 C.int) *C.go_perl_sv {
	cstr := C.CBytes(val)
	defer C.free(unsafe.Pointer(cstr))
	return C.go_perl_new_mortal_sv_string(gpi.pi, (*C.char)(cstr), C.go_perl_strlen(len(val)), wantUtf8)
}

func toPerlArgScalar(gpi *Interpreter, goValue interface{}) (*C.go_perl_sv, error) {
	switch val := goValue.(type) {
	case int:
		return C.go_perl_new_mortal_sv_iv(gpi.pi, C.IV(val)), nil
	case int16:
		return C.go_perl_new_mortal_sv_iv(gpi.pi, C.IV(val)), nil
	case int32:
		return C.go_perl_new_mortal_sv_iv(gpi.pi, C.IV(val)), nil
	case int64:
		return C.go_perl_new_mortal_sv_iv(gpi.pi, C.IV(val)), nil
	case uint:
		return C.go_perl_new_mortal_sv_uv(gpi.pi, C.UV(val)), nil
	case uint16:
		return C.go_perl_new_mortal_sv_uv(gpi.pi, C.UV(val)), nil
	case uint32:
		return C.go_perl_new_mortal_sv_uv(gpi.pi, C.UV(val)), nil
	case uint64:
		return C.go_perl_new_mortal_sv_uv(gpi.pi, C.UV(val)), nil
	case float32:
		return C.go_perl_new_mortal_sv_nv(gpi.pi, C.NV(val)), nil
	case float64:
		return C.go_perl_new_mortal_sv_nv(gpi.pi, C.NV(val)), nil
	case []byte:
		return toPerlMortalBinary(gpi, val, 0), nil
	case string:
		return toPerlMortalString(gpi, val, 1), nil
	case stringCoercion:
		var sv *C.go_perl_sv
		if val.bytes != nil {
			sv = toPerlMortalBinary(gpi, val.bytes, val.wantUtf8)
		} else {
			sv = toPerlMortalString(gpi, val.str, val.wantUtf8)
		}
		if sv == nil {
			return nil, ErrInvalidUtf8String
		}
		return sv, nil
	default:
		ty := reflect.TypeOf(goValue)
		name := ty.Name()
		if pkg := ty.PkgPath(); pkg != "" {
			return nil, fmt.Errorf("Unable to map %s.%s to a Perl scalar", pkg, name)
		}
		return nil, fmt.Errorf("Unable to map %s to a Perl scalar", name)
	}
}

func toPerlArgList(gpi *Interpreter, goValues []interface{}) (**C.go_perl_sv, error) {
	if len(goValues) == 0 {
		return nil, nil
	}

	memory := C.malloc(C.size_t(C.sizeof_go_perl_pointer_type * (len(goValues) + 1)))
	perlValues := (**C.go_perl_sv)(memory)
	for i, goValue := range goValues {
		var err error
		perlValue := (**C.go_perl_sv)(unsafe.Pointer(uintptr(memory) + uintptr(C.sizeof_go_perl_pointer_type*i)))
		*perlValue, err = toPerlArgScalar(gpi, goValue)
		if err != nil {
			return nil, err
		}
	}
	return perlValues, nil
}

func newScalarSliceFromMortals(gpi *Interpreter, count int, values **C.go_perl_sv) []*Scalar {
	slice := make([]*Scalar, count)
	for i := 0; i < count; i++ {
		result := (**C.go_perl_sv)(unsafe.Pointer(uintptr(unsafe.Pointer(values)) + uintptr(C.sizeof_go_perl_pointer_type*i)))
		slice[i] = newScalarFromMortal(gpi, *result)
	}
	return slice
}
