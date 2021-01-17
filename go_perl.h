#ifndef _GO_PERL_H_DEFINED
#define _GO_PERL_H_DEFINED

#include "EXTERN.h"
#include "perl.h"
#include "XSUB.h"

typedef PerlInterpreter go_perl_interpreter;
typedef SV go_perl_sv;
typedef STRLEN go_perl_strlen;

void go_perl_init();
go_perl_interpreter *go_perl_new_interpreter();

go_perl_sv *go_perl_get_global_scalar(go_perl_interpreter* my_perl, const char *name, int flags);

char *go_perl_svpv(go_perl_interpreter* my_perl, go_perl_sv* sv, go_perl_strlen* len);

#endif
