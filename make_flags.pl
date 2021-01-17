#!/usr/bin/env perl

use strict;
use warnings;

use FindBin qw();
use ExtUtils::Embed;

my $version = $^V;
my $cflags = ' ' . ExtUtils::Embed::ccopts . ' ';
my $lflags = ' ' . ExtUtils::Embed::ldopts . ' ';

$cflags =~ s/ -f\S+ / /g;
$lflags =~ s/ -[fW]\S+ / /g;

my $out_path = $FindBin::RealDir . '/go_perl_flags.go';
open my $out, '>', $out_path or die "Error opening $out_path: $!";

print $out <<EOT
package perl

/*
#cgo CFLAGS:  $cflags
#cgo LDFLAGS: $lflags
*/
import "C"
EOT
