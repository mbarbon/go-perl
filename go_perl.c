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

go_perl_sv* go_perl_new_mortal_sv_string(go_perl_interpreter* my_perl, char* pv, STRLEN len) {
    return sv_2mortal(newSVpvn(pv, len));
}

char* go_perl_svpv(go_perl_interpreter* my_perl, go_perl_sv* sv, go_perl_strlen* len) {
    return SvPV(sv, *len);
}

int go_perl_eval_void(go_perl_interpreter* my_perl, go_perl_sv* code, go_perl_sv** exc) {
    I32 retcount = do_eval(my_perl, code, G_VOID, exc);
    if (retcount < 0) {
        return retcount;
    }

    return 0;
}

int go_perl_call_void(go_perl_interpreter* my_perl, go_perl_sv* sub, int argcount, go_perl_sv** args, go_perl_sv** exc) {
    I32 retcount = do_call(my_perl, sub, G_VOID, argcount, args, exc);
    if (retcount < 0) {
        return retcount;
    }

    return 0;
}
