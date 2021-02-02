package perl

/*
#include "go_perl.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// Interpreter is an instance of the Perl interpreter
type Interpreter struct {
	pi *C.go_perl_interpreter
}

// NewInterpreter creates a new Perl interpreter
func NewInterpreter() *Interpreter {
	pi := C.go_perl_new_interpreter()
	gpi := &Interpreter{
		pi: pi,
	}
	runtime.SetFinalizer(gpi, func(gpi *Interpreter) {
		C.perl_destruct(gpi.pi)
		C.perl_free(gpi.pi)
	})
	return gpi
}

// Scalar returns a reference to an existing package variable, or nil
func (gpi *Interpreter) Scalar(name string) *Scalar {
	return gpi.scalar(name, 0)
}

// CreateScalar returns a reference to an existing package variable, or creates a new one
func (gpi *Interpreter) CreateScalar(name string) *Scalar {
	return gpi.scalar(name, C.GV_ADD)
}

func (gpi *Interpreter) scalar(name string, flags int) *Scalar {
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))
	sv := C.go_perl_get_global_scalar(gpi.pi, cstr, C.int(flags))
	if sv == nil {
		return nil
	}
	return newScalar(gpi, sv)
}
