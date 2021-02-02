package perl

/*
#include "go_perl.h"
*/
import "C"
import "unsafe"

func toPerlMortalString(gpi *Interpreter, val string) *C.go_perl_sv {
	cstr := C.CString(val)
	defer C.free(unsafe.Pointer(cstr))
	return C.go_perl_new_mortal_sv_string(gpi.pi, cstr, C.go_perl_strlen(len(val)))
}
