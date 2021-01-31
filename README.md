# go-perl
A silly experiment for running Perl from Go.

# Build
`go-perl` can determine compiler/linker flags for Perl in two ways: onet more convenient
when developing `go-perl` itself, the other more suitable for using `go-perl`.

## Build with a generated flag file
Generate a file named `go_perl_flags.go` containing Perl build flags.
```
$ perl make_flags.pl --go-perl-flags
```
Build with the `goperlflags` tag.
```
$ go build -tags goperlflags .
```
## Build with pkg-config flags
Generate a pkg-config `.pc` file for the current Perl.
```
$ perl make_flags.pl --pkg-config --pkg-config-path=/pkgconfig/file/path
```
Use the `PKG_CONFIG_PATH` environemnt variable so `go` can find the `.pc` file.
```
$ export PKG_CONFIG_PATH=/pkgconfig/file/path
$ go build .
```