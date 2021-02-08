package perl

/*
#include "go_perl.h"
*/
import "C"
import "runtime"

// Array wraps a Perl array value
type Array struct {
	gpi *Interpreter
	av  *C.go_perl_av
}

func newArray(gpi *Interpreter, av *C.go_perl_av) *Array {
	return &Array{
		gpi: gpi,
		av:  av,
	}
}

func newArrayFromMortal(gpi *Interpreter, av *C.go_perl_av) *Array {
	a := newArray(gpi, av)
	C.go_perl_sv_refcnt_inc(gpi.pi, C.av_to_sv(av))
	a.setFinalizer()
	return a
}

// Fetch returns the index-th element of the array
func (a *Array) Fetch(index int) *Scalar {
	valPtr := C.go_perl_av_fetch(a.gpi.pi, a.av, C.I32(index))
	if valPtr == nil {
		return nil
	}
	return newScalarFromMortal(a.gpi, *valPtr)
}

func (a *Array) setFinalizer() {
	runtime.SetFinalizer(a, func(a *Array) {
		C.go_perl_sv_refcnt_dec(a.gpi.pi, C.av_to_sv(a.av))
	})
}
