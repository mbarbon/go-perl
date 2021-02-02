package perl

/*
#include "go_perl.h"
*/
import "C"

func init() {
	C.go_perl_init()
}
