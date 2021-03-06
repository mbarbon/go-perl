#ifndef _GO_PERL_H_DEFINED
#define _GO_PERL_H_DEFINED

#define PERL_NO_GET_CONTEXT
#include "EXTERN.h"
#include "perl.h"
#include "XSUB.h"

#define GO_PERL_SVf_IV     0x0001
#define GO_PERL_SVf_UV     0x0002
#define GO_PERL_SVf_NV     0x0004
#define GO_PERL_SVf_PV     0x0008
#define GO_PERL_SVf_PVutf8 0x0010
#define GO_PERL_SVf_RV     0x0020

typedef int* go_perl_pointer_type;
typedef PerlInterpreter go_perl_interpreter;
typedef SV go_perl_sv;
typedef AV go_perl_av;
typedef HV go_perl_hv;
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

int go_perl_sv_type(go_perl_interpreter* my_perl, go_perl_sv* sv);
int go_perl_sv_flags(go_perl_interpreter* my_perl, go_perl_sv* sv);
int go_perl_decoded_sv_flags(go_perl_interpreter* my_perl, go_perl_sv* sv);

go_perl_sv* go_perl_new_mortal_sv_iv(go_perl_interpreter* my_perl, IV iv);
go_perl_sv* go_perl_new_mortal_sv_uv(go_perl_interpreter* my_perl, UV iv);
go_perl_sv* go_perl_new_mortal_sv_nv(go_perl_interpreter* my_perl, NV iv);
go_perl_sv* go_perl_new_mortal_sv_string(go_perl_interpreter* my_perl, char* pv, STRLEN len, int want_utf8);

go_perl_sv* go_perl_svrv(go_perl_interpreter* my_perl, go_perl_sv* sv);
IV go_perl_sviv(go_perl_interpreter* my_perl, go_perl_sv* sv);
UV go_perl_svuv(go_perl_interpreter* my_perl, go_perl_sv* sv);
NV go_perl_svnv(go_perl_interpreter* my_perl, go_perl_sv* sv);
char* go_perl_svpv(go_perl_interpreter* my_perl, go_perl_sv* sv, go_perl_strlen* len);
char* go_perl_svpv_bytes(go_perl_interpreter* my_perl, go_perl_sv* sv, go_perl_strlen* len);
char* go_perl_svpv_utf8(go_perl_interpreter* my_perl, go_perl_sv* sv, go_perl_strlen* len);

go_perl_sv** go_perl_av_fetch(go_perl_interpreter* my_perl, go_perl_av* av, I32 index);

go_perl_sv** go_perl_hv_fetch(go_perl_interpreter* my_perl, go_perl_hv* hv, const char *key, I32 klen);

int go_perl_eval_void(go_perl_interpreter* my_perl, go_perl_sv* code, go_perl_sv** exc);
int go_perl_eval_scalar(go_perl_interpreter* my_perl, go_perl_sv* code, go_perl_sv** exc, go_perl_sv** result);
int go_perl_eval_list(go_perl_interpreter* my_perl, go_perl_sv* code, go_perl_sv** exc, go_perl_sv*** results);

int go_perl_call_void(go_perl_interpreter* my_perl, go_perl_sv* sub, int argcount, go_perl_sv** args, go_perl_sv** exc);
int go_perl_call_scalar(go_perl_interpreter* my_perl, go_perl_sv* sub, int argcount, go_perl_sv** args, go_perl_sv** exc, go_perl_sv** result);
int go_perl_call_list(go_perl_interpreter* my_perl, go_perl_sv* sub, int argcount, go_perl_sv** args, go_perl_sv** exc, go_perl_sv*** results);

// casts

inline go_perl_sv* cv_to_sv(go_perl_cv* cv) {
    return (go_perl_sv*) cv;
}

inline go_perl_cv* sv_to_cv(go_perl_sv* sv) {
    return (go_perl_cv*) sv;
}

inline go_perl_sv* av_to_sv(go_perl_av* av) {
    return (go_perl_sv*) av;
}

inline go_perl_av* sv_to_av(go_perl_sv* sv) {
    return (go_perl_av*) sv;
}

inline go_perl_sv* hv_to_sv(go_perl_hv* hv) {
    return (go_perl_sv*) hv;
}

inline go_perl_hv* sv_to_hv(go_perl_sv* sv) {
    return (go_perl_hv*) sv;
}

#endif
