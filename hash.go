package perl

/*
#include "go_perl.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// Hash wraps a Perl hash value
type Hash struct {
	gpi *Interpreter
	hv  *C.go_perl_hv
}

func newHash(gpi *Interpreter, hv *C.go_perl_hv) *Hash {
	return &Hash{
		gpi: gpi,
		hv:  hv,
	}
}

func newHashFromMortal(gpi *Interpreter, hv *C.go_perl_hv) *Hash {
	h := newHash(gpi, hv)
	C.go_perl_sv_refcnt_inc(gpi.pi, C.hv_to_sv(hv))
	h.setFinalizer()
	return h
}

// FetchStringKey fetches the given key from the hash
func (h *Hash) FetchStringKey(key string) *Scalar {
	cstr := C.CString(key)
	defer C.free(unsafe.Pointer(cstr))
	cstrLen := C.I32(len(key))
	if C.go_perl_looks_like_utf8(cstr, C.STRLEN(cstrLen)) {
		cstrLen = -cstrLen
	}
	valPtr := C.go_perl_hv_fetch(h.gpi.pi, h.hv, cstr, cstrLen)
	if valPtr == nil {
		return nil
	}
	return newScalarFromMortal(h.gpi, *valPtr)
}

func (h *Hash) setFinalizer() {
	runtime.SetFinalizer(h, func(h *Hash) {
		C.go_perl_sv_refcnt_dec(h.gpi.pi, C.hv_to_sv(h.hv))
	})
}
