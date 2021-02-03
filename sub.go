package perl

/*
#include "go_perl.h"
*/
import "C"
import "unsafe"

// Sub wraps a Perl subroutine value
type Sub struct {
	gpi *Interpreter
	cv  *C.go_perl_cv
}

func newSub(gpi *Interpreter, cv *C.go_perl_cv) *Sub {
	return &Sub{
		gpi: gpi,
		cv:  cv,
	}
}

// CallVoid calls this sub in void context
func (s *Sub) CallVoid(goArgs ...interface{}) error {
	C.go_perl_open_scope(s.gpi.pi)
	defer C.go_perl_close_scope(s.gpi.pi)

	perlArgs, err := toPerlArgList(s.gpi, goArgs)
	if err != nil {
		return err
	}
	defer C.free(unsafe.Pointer(perlArgs))

	var exc *C.go_perl_sv
	var rescount = C.go_perl_call_void(s.gpi.pi, (*C.go_perl_sv)(unsafe.Pointer(s.cv)), C.int(len(goArgs)), perlArgs, &exc)
	if rescount < 0 {
		return s.gpi.evalError(rescount, exc)
	}

	return nil
}
