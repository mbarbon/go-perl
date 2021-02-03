package perl

/*
#include "go_perl.h"
*/
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

// ErrExitCalled is returned from Eval*/Call* when the Perl code call exit()
var ErrExitCalled = errors.New("exit called")

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

// Sub returns a reference to an existing package variable, or nil
func (gpi *Interpreter) Sub(name string) *Sub {
	return gpi.sub(name, 0)
}

// CreateSub returns a reference to an existing package variable, or creates a new one
func (gpi *Interpreter) CreateSub(name string) *Sub {
	return gpi.sub(name, C.GV_ADD)
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

func (gpi *Interpreter) sub(name string, flags int) *Sub {
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))
	cv := C.go_perl_get_global_code(gpi.pi, cstr, C.int(flags))
	if cv == nil {
		return nil
	}
	return newSub(gpi, cv)
}

// EvalVoid avals the given string as Perl code, in void context
func (gpi *Interpreter) EvalVoid(code string) error {
	C.go_perl_open_scope(gpi.pi)
	defer C.go_perl_close_scope(gpi.pi)

	perlCode := toPerlMortalString(gpi, code)

	var exc *C.go_perl_sv
	var rescount = C.go_perl_eval_void(gpi.pi, perlCode, &exc)
	if rescount < 0 {
		return gpi.evalError(rescount, exc)
	}

	return nil
}

func (gpi *Interpreter) evalError(errCode C.int, exc *C.go_perl_sv) error {
	if errCode == -2 {
		return ErrExitCalled
	}
	return newScalarFromMortal(gpi, exc).asError()
}
