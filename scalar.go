package perl

/*
#include "go_perl.h"
*/
import "C"
import "runtime"

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
	str := C.go_perl_svpv(s.gpi.pi, s.sv, &len)
	return C.GoStringN(str, C.int(len))
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
