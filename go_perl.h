#ifndef _GO_PERL_H_DEFINED
#define _GO_PERL_H_DEFINED

#define PERL_NO_GET_CONTEXT
#include "EXTERN.h"
#include "perl.h"
#include "XSUB.h"

typedef PerlInterpreter go_perl_interpreter;
typedef SV go_perl_sv;
typedef STRLEN go_perl_strlen;

void go_perl_init();
go_perl_interpreter* go_perl_new_interpreter();

void go_perl_open_scope(go_perl_interpreter* my_perl);
void go_perl_close_scope(go_perl_interpreter* my_perl);

go_perl_sv* go_perl_get_global_scalar(go_perl_interpreter* my_perl, const char* name, int flags);

void go_perl_sv_refcnt_inc(go_perl_interpreter* my_perl, go_perl_sv* sv);
void go_perl_sv_refcnt_dec(go_perl_interpreter* my_perl, go_perl_sv* sv);

go_perl_sv* go_perl_new_mortal_sv_string(go_perl_interpreter* my_perl, char* pv, STRLEN len);

char* go_perl_svpv(go_perl_interpreter* my_perl, go_perl_sv* sv, go_perl_strlen* len);

int go_perl_eval_void(go_perl_interpreter* my_perl, go_perl_sv* code, go_perl_sv** exc);

#endif
