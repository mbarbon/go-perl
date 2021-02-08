#include "go_perl.h"

#define EVAL_DIE -1
#define EVAL_EXIT -2

static int go_perl_argc = 3;
static char* go_perl_argv_arr[4] = { "go-perl", "-e", "0", NULL};
static char* go_perl_env_arr[1] = { NULL };

static char** go_perl_argv = go_perl_argv_arr;
static char** go_perl_env = go_perl_env_arr;

static I32 do_eval(go_perl_interpreter* my_perl, go_perl_sv* code, I32 flags, go_perl_sv** exc) {
    int ret;
    dJMPENV;

    JMPENV_PUSH(ret);

    I32 retcount;
    switch (ret) {
    case 0:
        retcount = eval_sv(code, flags);
        go_perl_sv* errsv = GvSVn(PL_errgv);
        if (!SvOK(errsv) || !SvPOK(errsv) || SvCUR(errsv) > 0) {
            *exc = errsv;
            retcount = EVAL_DIE;
        }
        break;
    case 1:
        // complete failure
    case 2:
        // exit called
        retcount = EVAL_EXIT;
        break;
    case 3:
        // this should only happen when G_RETHROW is used
        *exc = GvSV(PL_errgv);
        retcount = EVAL_DIE;
        break;
    }

    JMPENV_POP;

    return retcount;
}

static I32 do_call(go_perl_interpreter* my_perl, go_perl_sv* sub, I32 flags, int argcount, go_perl_sv** args, go_perl_sv**exc) {
    int ret;
    dJMPENV;
    dSP;

    PUSHMARK(SP);
    EXTEND(SP, argcount);
    for (int i = 0; i < argcount; ++i) {
        PUSHs(args[i]);
    }
    PUTBACK;

    JMPENV_PUSH(ret);

    I32 retcount;
    switch (ret) {
    case 0:
        retcount = call_sv(sub, flags|G_EVAL);
        go_perl_sv* errsv = GvSVn(PL_errgv);
        if (!SvOK(errsv) || !SvPOK(errsv) || SvCUR(errsv) > 0) {
            *exc = errsv;
            retcount = EVAL_DIE;
        }
        break;
    case 1:
        // complete failure
    case 2:
        // exit called
        retcount = EVAL_EXIT;
        break;
    case 3:
        // this should never happen
        *exc = GvSV(PL_errgv);
        retcount = EVAL_DIE;
        break;
    }

    JMPENV_POP;

    return retcount;
}

static int return_list(go_perl_interpreter* my_perl, I32 retcount, go_perl_sv*** result) {
    if (retcount != 0) {
        dSP;
        go_perl_sv** res = (go_perl_sv**) malloc(sizeof(go_perl_sv*) * retcount);
        for (int i = retcount - 1; i >= 0; --i) {
            res[i] = POPs;
        }
        *result = res;
        PUTBACK;
    }

    return retcount;
}

#ifndef is_invariant_string_log

static bool is_invariant_string_loc(const U8* s, STRLEN len, const U8** ep) {
    const U8* send = s + len;

    while (s < send && UTF8_IS_INVARIANT(*s)) {
        ++s;
    }
    *ep = s;

    return send == s;
}

#endif

// 2: UTF-8 with high-bit chars
// 1: UTF-8 invariant chars (< 128, could be ASCII)
// 0: non-UTF-8, non-ASCII; could be Latin-1 or binary, no way to know
static int utf8_kind(U8* s, STRLEN len) {
    const U8* end;

    if (is_invariant_string_loc(s, len, &end)) {
        return 1;
    }
    if (is_utf8_string(end, len - (end - s))) {
        return 2;
    }

    return 0;
}

void go_perl_init() {
    PERL_SYS_INIT3(&go_perl_argc, &go_perl_argv, &go_perl_env);
}

go_perl_interpreter* go_perl_new_interpreter() {
    go_perl_interpreter* my_perl;

    my_perl = perl_alloc();
    perl_construct(my_perl);
    PL_exit_flags |= PERL_EXIT_DESTRUCT_END;
    perl_parse(my_perl, NULL, go_perl_argc, go_perl_argv, (char**)NULL);
    perl_run(my_perl);

    return my_perl;
}

void go_perl_open_scope(go_perl_interpreter* my_perl) {
    ENTER;
    SAVETMPS;
}

void go_perl_close_scope(go_perl_interpreter* my_perl) {
    FREETMPS;
    LEAVE;
}

go_perl_sv* go_perl_get_global_scalar(go_perl_interpreter* my_perl, const char* name, int flags) {
    return get_sv(name, flags);
}

go_perl_cv* go_perl_get_global_code(go_perl_interpreter* my_perl, const char* name, int flags) {
    return get_cv(name, flags);
}

void go_perl_sv_refcnt_inc(go_perl_interpreter* my_perl, go_perl_sv* sv) {
    SvREFCNT_inc(sv);
}

void go_perl_sv_refcnt_dec(go_perl_interpreter* my_perl, go_perl_sv* sv) {
    SvREFCNT_dec(sv);
}

int go_perl_sv_type(go_perl_interpreter* my_perl, go_perl_sv* sv) {
    return SvTYPE(sv);
}

int go_perl_sv_flags(go_perl_interpreter* my_perl, go_perl_sv* sv) {
    return SvFLAGS(sv);
}

int go_perl_decoded_sv_flags(go_perl_interpreter* my_perl, go_perl_sv* sv) {
    if (SvROK(sv)) {
        return GO_PERL_SVf_RV;
    }

    int decoded_flags = 0;

    if (SvNOK(sv)) {
        decoded_flags |= GO_PERL_SVf_NV;
    }
    if (SvIOK(sv)) {
        decoded_flags |= GO_PERL_SVf_IV;
    }
    if (SvUOK(sv)) {
        decoded_flags |= GO_PERL_SVf_UV;
    }
    if (SvPOK(sv)) {
        decoded_flags |= GO_PERL_SVf_PV;
        if (SvUTF8(sv)) {
        decoded_flags |= GO_PERL_SVf_PVutf8;
        }
    }

    return decoded_flags;
}

go_perl_sv* go_perl_new_mortal_sv_iv(go_perl_interpreter* my_perl, IV iv){
    return sv_2mortal(newSViv(iv));
}

go_perl_sv* go_perl_new_mortal_sv_uv(go_perl_interpreter* my_perl, UV uv){
    return sv_2mortal(newSVuv(uv));
}

go_perl_sv* go_perl_new_mortal_sv_nv(go_perl_interpreter* my_perl, NV nv){
    return sv_2mortal(newSVnv(nv));
}

go_perl_sv* go_perl_new_mortal_sv_string(go_perl_interpreter* my_perl, char* pv, STRLEN len, int want_utf8) {
    int have_utf8 = utf8_kind(pv, len);
    if (want_utf8 == 2 && have_utf8 == 0) {
        return NULL;
    }

    go_perl_sv* sv = sv_2mortal(newSVpvn(pv, len));
    if (want_utf8) {
        if ((want_utf8 == 2 && have_utf8 > 0) || have_utf8 == 2) {
            SvUTF8_on(sv);
        } else if (want_utf8 == 2) {
            return NULL;
        }
    }

    return sv;
}

char* go_perl_svpv(go_perl_interpreter* my_perl, go_perl_sv* sv, go_perl_strlen* len) {
    return SvPV(sv, *len);
}

char* go_perl_svpv_bytes(go_perl_interpreter* my_perl, go_perl_sv* sv, go_perl_strlen* len) {
    return SvPVbyte(sv, *len);
}

char* go_perl_svpv_utf8(go_perl_interpreter* my_perl, go_perl_sv* sv, go_perl_strlen* len) {
    return SvPVutf8(sv, *len);
}

go_perl_sv* go_perl_svrv(go_perl_interpreter* my_perl, go_perl_sv* sv) {
    return SvRV(sv);
}

IV go_perl_sviv(go_perl_interpreter* my_perl, go_perl_sv* sv) {
    return SvIV(sv);
}

UV go_perl_svuv(go_perl_interpreter* my_perl, go_perl_sv* sv) {
    return SvUV(sv);
}

NV go_perl_svnv(go_perl_interpreter* my_perl, go_perl_sv* sv) {
    return SvNV(sv);
}

go_perl_sv** go_perl_av_fetch(go_perl_interpreter* my_perl, go_perl_av* av, I32 index) {
    return av_fetch(av, index, 0);
}

go_perl_sv** go_perl_hv_fetch(go_perl_interpreter* my_perl, go_perl_hv* hv, const char *key, I32 klen) {
    return hv_fetch(hv, key, klen, 0);
}

int go_perl_eval_void(go_perl_interpreter* my_perl, go_perl_sv* code, go_perl_sv** exc) {
    I32 retcount = do_eval(my_perl, code, G_VOID, exc);
    if (retcount < 0) {
        return retcount;
    }

    return 0;
}

int go_perl_eval_scalar(go_perl_interpreter* my_perl, go_perl_sv* code, go_perl_sv**exc, go_perl_sv** result) {
    I32 retcount = do_eval(my_perl, code, G_SCALAR, exc);
    if (retcount < 0) {
        return retcount;
    }

    dSP;
    *result = POPs;
    PUTBACK;

    return 0;
}

int go_perl_eval_list(go_perl_interpreter* my_perl, go_perl_sv* code, go_perl_sv** exc, go_perl_sv*** result) {
    I32 retcount = do_eval(my_perl, code, G_ARRAY, exc);
    if (retcount < 0) {
        return retcount;
    }

    return return_list(my_perl, retcount, result);
}

int go_perl_call_void(go_perl_interpreter* my_perl, go_perl_sv* sub, int argcount, go_perl_sv** args, go_perl_sv** exc) {
    I32 retcount = do_call(my_perl, sub, G_VOID, argcount, args, exc);
    if (retcount < 0) {
        return retcount;
    }

    return 0;
}

int go_perl_call_scalar(go_perl_interpreter* my_perl, go_perl_sv* sub, int argcount, go_perl_sv** args, go_perl_sv**exc, go_perl_sv** result) {
    I32 retcount = do_call(my_perl, sub, G_SCALAR, argcount, args, exc);
    if (retcount < 0) {
        return retcount;
    }

    dSP;
    *result = POPs;
    PUTBACK;

    return 0;
}

int go_perl_call_list(go_perl_interpreter* my_perl, go_perl_sv* sub, int argcount, go_perl_sv** args, go_perl_sv** exc, go_perl_sv*** result) {
    I32 retcount = do_call(my_perl, sub, G_ARRAY, argcount, args, exc);
    if (retcount < 0) {
        return retcount;
    }

    return return_list(my_perl, retcount, result);
}
