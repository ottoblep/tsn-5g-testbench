#!/bin/bash
#
# Configure iptables in UPF
#
# Forward everything
iptables -I FORWARD -j ACCEPT

# Allow return traffic
iptables -I INPUT -m state --state RELATED,ESTABLISHED -j ACCEPT