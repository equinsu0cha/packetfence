package pfconfig::namespaces::config::OrderedRealm;

=head1 NAME

pfconfig::namespaces::config::OrderedRealm

=cut

=head1 DESCRIPTION

pfconfig::namespaces::config::OrderedRealm

This module creates the configuration array associated to realm.conf

=cut

use strict;
use warnings;

use pfconfig::namespaces::config;
use pf::file_paths qw(
  $realm_default_config_file
  $realm_config_file
);

use base 'pfconfig::namespaces::config';

sub init {
    my ($self) = @_;
    $self->{_scoped_by_tenant_id} = 1;
    $self->{child_resources} = [ "resource::RealmReverseLookup" ];
    $self->{ini} = pf::IniFiles->new(
        -file       => $realm_config_file,
        -import     => pf::IniFiles->new(-file => $realm_default_config_file),
        -allowempty => 1,
    );
}

sub build {
    my ($self) = @_;
    return {
        map {
            my $id = $_;
            $id => $self->build_tenant_sections($id)
        } $self->tenant_ids()
    };
}

sub tenant_ids {
    my ($self) = @_;
    return $self->{ini}->Groups();
}

sub build_tenant_sections {
    my ($self, $tenant_id) = @_;
    return [
        map {
            my $section = $_;
            my $name = $section;
            $name =~ s/^\Q$tenant_id \E//;
            $name = lc($name);
        } $self->sections_for_tenant($tenant_id)
    ];
}

sub sections_for_tenant {
    my ($self, $tenant_id) = @_;
    return $self->{ini}->GroupMembers($tenant_id);
}

=head1 AUTHOR

Inverse inc. <info@inverse.ca>

=head1 COPYRIGHT

Copyright (C) 2005-2021 Inverse inc.

=head1 LICENSE

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301,
USA.

=cut

1;

# vim: set shiftwidth=4:
# vim: set expandtab:
# vim: set backspace=indent,eol,start:
