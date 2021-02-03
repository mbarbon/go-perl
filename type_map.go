package perl

/*
#include "go_perl.h"
*/
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

func toPerlMortalString(gpi *Interpreter, val string) *C.go_perl_sv {
	cstr := C.CString(val)
	defer C.free(unsafe.Pointer(cstr))
	return C.go_perl_new_mortal_sv_string(gpi.pi, cstr, C.go_perl_strlen(len(val)))
}

func toPerlArgScalar(gpi *Interpreter, goValue interface{}) (*C.go_perl_sv, error) {
	switch val := goValue.(type) {
	case string:
		return toPerlMortalString(gpi, val), nil
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
