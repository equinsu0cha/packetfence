#/bin/perl


my %snort_rules_version = (
    "RHEL5" => "snort-2.8.6",
    "RHEL6" => "snort-2.9.0",
    "latest" => "snort-2.9.0",
    "SQUEEZE" => "snort-2.8.4",
);

my %oses = (
    "CentOS release 5" => "RHEL5",
    "Red Hat Enterprise Linux Server release 5" => "RHEL5",
    "CentOS Linux release 6" => "RHEL6",
    "CentOS release 6" => "RHEL6",
    "Red Hat Enterprise Linux Server release 6" => "RHEL6",
    "6" => "SQUEEZE",
);

my $os_type = supported_os();

   my @rule_files = (
        'emerging-botcc.rules',
        'emerging-attack_response.rules',
        'emerging-exploit.rules',
        'emerging-malware.rules',
        'emerging-p2p.rules',
        'emerging-scan.rules',
        'emerging-shellcode.rules',
        'emerging-trojan.rules',
        'emerging-virus.rules',
        'emerging-worm.rules'
    );
    foreach my $current_rule_file (@rule_files) {
        `/usr/bin/wget -N http://rules.emergingthreats.net/open/$snort_rules_version{$os_type}/rules/$current_rule_file -P /usr/local/pf/conf/snort`;
    }

sub supported_os {

    # RedHat and derivatives
    if ( -e "/etc/redhat-release" ) {
        my $rhrelease_fh;
        open( $rhrelease_fh, '<', "/etc/redhat-release" );
        $version = <$rhrelease_fh>;
        close($rhrelease_fh);
    }
    # Debian and derivatives
    elsif (-e "/etc/debian_version" ) {
        my $debianversion;
        open( $debianversion, '<', "/etc/debian_version" );
        $version = <$debianversion>;
        close($debianversion);
    }
    # Unknown
    else {
        $version = "X";
    }

    foreach my $supported ( keys(%oses) ) {
        return ( $oses{$supported} ) if ( $version =~ /^$supported/ );
    }
    return (0);
}
