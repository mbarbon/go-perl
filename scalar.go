package perl

/*
#include "go_perl.h"
*/
import "C"

// Scalar wraps a Perl scalar value
type Scalar struct {
	gpi *Interpreter
	sv  *C.go_perl_sv
}

func newScalar(gpi *Interpreter, sv *C.go_perl_sv) *Scalar {
	return &Scalar{
		gpi: gpi,
		sv:  sv,
	}
}

// String returns the scalar value coerced to a string
func (s *Scalar) String() string {
	var len C.go_perl_strlen
	str := C.go_perl_svpv(s.gpi.pi, s.sv, &len)
	return C.GoStringN(str, C.int(len))
}
