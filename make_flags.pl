#!/usr/bin/env perl

use strict;
use warnings;

use FindBin qw();
use File::Path qw();
use ExtUtils::Embed;
use Getopt::Long;

sub help {
    my ($code) = @_;

    print <<EOT;
Usage: $0 <--go-perl-flags|--pkg-config> [OPTIONS]

Write compiler/linker flags to build against Perl.

  --help
      Display this help message and exit.

  --go-perl-flags
      Write build flags for the current Perl versions to go_perl_flags.go

  --pkg-config
      Write a pkg-config .pc file with build flags for the current
      Perl version to the .pkgconfig directory.

  --pkg-config-path=PATH
      Change the directory where --pkg-config writes the .pc file.

EOT

    exit $code;
}

sub main {
    GetOptions(
        'help'              => \my $help,
        'pkg-config'        => \my $pkg_config,
        'pkg-config-path=s' => \my $pkg_config_path,
        'go-perl-flags'     => \my $go_perl_flags,
    ) or help(1);

    if ($help) {
        help(0);
    }
    if (!$pkg_config && !$go_perl_flags) {
        help(1);
    }
    if ($go_perl_flags) {
        write_go_perl_flags();
    }
    if ($pkg_config) {
        write_pkg_config($pkg_config_path);
    }
}

sub get_flags {
    my $version = $^V;
    my $cflags = ' ' . ExtUtils::Embed::ccopts . ' ';
    my $lflags = ' ' . ExtUtils::Embed::ldopts . ' ';

    $cflags =~ s/ -f\S+ / /g;
    $lflags =~ s/ -[fW]\S+ / /g;
    $cflags =~ s/^\s+|\s+$//g;
    $lflags =~ s/^\s+|\s+$//g;

    return ($version, $cflags, $lflags);
}

sub write_go_perl_flags {
    my (undef, $cflags, $lflags) = get_flags();
    my $out_path = $FindBin::RealDir . '/go_perl_flags.go';
    open my $out, '>', $out_path or die "Error opening $out_path: $!";

    print $out <<EOT
// +build goperlflags

package perl

/*
#cgo CFLAGS:  $cflags
#cgo LDFLAGS: $lflags
*/
import "C"
EOT
}

sub write_pkg_config {
    my ($pkg_config_path) = @_;

    if (!$pkg_config_path) {
        $pkg_config_path = $FindBin::RealDir . '/.pkgconfig';
        File::Path::make_path($pkg_config_path);
    }

    my ($version, $cflags, $lflags) = get_flags();
    my $out_path = $pkg_config_path . '//go-perl-perl.pc';
    open my $out, '>', $out_path or die "Error opening $out_path: $!";

    print $out <<EOT
Name: go-perl-perl
Description: Perl build flags for the go-perl package
Version: $version
Libs: $lflags
Cflags: $cflags
EOT
}

main();
exit 0;
