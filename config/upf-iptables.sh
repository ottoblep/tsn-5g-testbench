#!/bin/bash
#
# Configure iptables in UPF
#
# Forward everything
iptables -A FORWARD -j ACCEPT

# Masquerade outgoing traffic
iptables -t nat -A POSTROUTING -o eth1 -j MASQUERADE
# iptables -t nat -A POSTROUTING -o upfgtp -j MASQUERADE

# Allow return traffic
iptables -A INPUT -i upfgtp -m state --state RELATED,ESTABLISHED -j ACCEPT
iptables -A INPUT -i eth1 -m state --state RELATED,ESTABLISHED -j ACCEPT