package pfappserver::Form::User;

=head1 NAME

pfappserver::Form::User - Web form for a user

=head1 DESCRIPTION

Form definition to update a user.

=cut

use pf::config;
use pf::util qw(get_abbr_time);
use HTTP::Status qw(:constants is_success);
use HTML::FormHandler::Moose;
extends 'pfappserver::Base::Form::Authentication::Action';
has '+source_type' => ( default => 'SQL' );

=head2 Fields

=over

=item pid
=cut

has_field 'pid' =>
  (
   type => 'Uneditable',
   label => 'Username',
  );

=item firstname
=cut

has_field 'firstname' =>
  (
   type => 'Text',
   label => 'Firstname',
  );

=item lastname
=cut

has_field 'lastname' =>
  (
   type => 'Text',
   label => 'Lastname',
  );

=item company
=cut

has_field 'company' =>
  (
   type => 'Text',
   label => 'Company',
  );

=item email
=cut

has_field 'email' =>
  (
   type => 'Email',
   label => 'Email',
   required => 1,
  );

=item sponsor
=cut

has_field 'sponsor' =>
  (
   type => 'Text',
   label => 'Sponsor',
  );

=item address
=cut

has_field 'address' =>
  (
   type => 'TextArea',
   label => 'Address',
  );

=item notes
=cut

has_field 'notes' =>
  (
   type => 'TextArea',
   label => 'Notes',
  );

=item valid_from
=cut

has_field 'valid_from' =>
  (
   type => 'DatePicker',
   label => 'Arrival Date',
   required => 1,
   start => &now,
  );

=back

=head2 Blocks
=cut

=over
=cut


=item user block
  The user block contains the static fields of user
=cut

has_block 'user' =>
  (
   render_list => [qw(pid firstname lastname company email sponsor address notes)],
  );

=item templates block

 The templates block contains the dynamic fields of the rule definition.

 The following fields depend on the selected condition attribute :
  - the condition operators select fields
  - the condition value fields
 The following fields depend on the selected action type :
  - the action value fields

=cut

has_block 'templates' =>
  (
   tag => 'div',
   render_list => [
                   map( { "${_}_action" } @Actions::ACTIONS), # the field are defined in the super class
                  ],
   attr => { id => 'templates' },
   class => [ 'hidden' ],
  );

=back

=head2 Methods

=over

=item options_access_level

Populate the select field for the 'access level' template action.

=cut

sub options_access_level {
    my $self = shift;

    return
      (
       {
        label => $self->_localize('None'),
        value => $WEB_ADMIN_NONE,
       },
       {
        label => $self->_localize('All'),
        value => $WEB_ADMIN_ALL,
       },
      );
}

=item options_roles

Populate the select field for the roles template action.

=cut

sub options_roles {
    my $self = shift;

    my @roles;

    # Build a list of existing roles
    my ($status, $result) = $self->form->ctx->model('Roles')->list();
    if (is_success($status)) {
        @roles = map { $_->{name} => $_->{name} } @$result;
    }

    return @roles;
}

=item options_durations

Populate the access duration select field with the available values defined
in the pf.conf configuration file.

=cut

sub options_durations {
    my $self = shift;

    my $durations = pf::web::util::get_translated_time_hash(
        [ split (/\s*,\s*/, $Config{'guests_admin_registration'}{'access_duration_choices'}) ],
        $self->form->ctx->languages()->[0]
    );
    my @options = map { get_abbr_time($_) => $durations->{$_} } sort { $a <=> $b } keys %$durations;

    return \@options;
}

=item now
Return the current day, used as the minimal date of the arrival date.
=cut

sub now {
    my ($sec,$min,$hour,$mday,$mon,$year) = localtime(time);
    return sprintf "%d-%02d-%02d", $year+1900, $mon+1, $mday;
}

=back

=head1 COPYRIGHT

Copyright (C) 2012 Inverse inc.

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

__PACKAGE__->meta->make_immutable;
1;
