#ifndef _GO_PERL_H_DEFINED
#define _GO_PERL_H_DEFINED

#define PERL_NO_GET_CONTEXT
#include "EXTERN.h"
#include "perl.h"
#include "XSUB.h"

typedef int* go_perl_pointer_type;
typedef PerlInterpreter go_perl_interpreter;
typedef SV go_perl_sv;
typedef CV go_perl_cv;
typedef STRLEN go_perl_strlen;

void go_perl_init();
go_perl_interpreter* go_perl_new_interpreter();

void go_perl_open_scope(go_perl_interpreter* my_perl);
void go_perl_close_scope(go_perl_interpreter* my_perl);

go_perl_sv* go_perl_get_global_scalar(go_perl_interpreter* my_perl, const char* name, int flags);
go_perl_cv* go_perl_get_global_code(go_perl_interpreter* my_perl, const char* name, int flags);

void go_perl_sv_refcnt_inc(go_perl_interpreter* my_perl, go_perl_sv* sv);
void go_perl_sv_refcnt_dec(go_perl_interpreter* my_perl, go_perl_sv* sv);

go_perl_sv* go_perl_new_mortal_sv_iv(go_perl_interpreter* my_perl, IV iv);
go_perl_sv* go_perl_new_mortal_sv_uv(go_perl_interpreter* my_perl, UV iv);
go_perl_sv* go_perl_new_mortal_sv_nv(go_perl_interpreter* my_perl, NV iv);
go_perl_sv* go_perl_new_mortal_sv_string(go_perl_interpreter* my_perl, char* pv, STRLEN len, int want_utf8);

IV go_perl_sviv(go_perl_interpreter* my_perl, go_perl_sv* sv);
UV go_perl_svuv(go_perl_interpreter* my_perl, go_perl_sv* sv);
NV go_perl_svnv(go_perl_interpreter* my_perl, go_perl_sv* sv);
char* go_perl_svpv(go_perl_interpreter* my_perl, go_perl_sv* sv, go_perl_strlen* len);
char* go_perl_svpv_bytes(go_perl_interpreter* my_perl, go_perl_sv* sv, go_perl_strlen* len);
char* go_perl_svpv_utf8(go_perl_interpreter* my_perl, go_perl_sv* sv, go_perl_strlen* len);

int go_perl_eval_void(go_perl_interpreter* my_perl, go_perl_sv* code, go_perl_sv** exc);
int go_perl_eval_scalar(go_perl_interpreter* my_perl, go_perl_sv* code, go_perl_sv** exc, go_perl_sv** result);
int go_perl_eval_list(go_perl_interpreter* my_perl, go_perl_sv* code, go_perl_sv** exc, go_perl_sv*** results);

int go_perl_call_void(go_perl_interpreter* my_perl, go_perl_sv* sub, int argcount, go_perl_sv** args, go_perl_sv** exc);
int go_perl_call_scalar(go_perl_interpreter* my_perl, go_perl_sv* sub, int argcount, go_perl_sv** args, go_perl_sv** exc, go_perl_sv** result);
int go_perl_call_list(go_perl_interpreter* my_perl, go_perl_sv* sub, int argcount, go_perl_sv** args, go_perl_sv** exc, go_perl_sv*** results);

#endif
