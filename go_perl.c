#include "go_perl.h"

static int go_perl_argc = 3;
static char* go_perl_argv_arr[4] = { "go-perl", "-e", "0", NULL};
static char* go_perl_env_arr[1] = { NULL };

static char** go_perl_argv = go_perl_argv_arr;
static char** go_perl_env = go_perl_env_arr;

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

go_perl_sv* go_perl_get_global_scalar(go_perl_interpreter* my_perl, const char* name, int flags) {
    return get_sv(name, flags);
}

char *go_perl_svpv(go_perl_interpreter* my_perl, go_perl_sv* sv, go_perl_strlen* len) {
    return SvPV(sv, *len);
}