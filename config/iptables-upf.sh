#!/bin/bash
#
# Configure iptables in UPF
#
# Forward everything
iptables -I FORWARD -j ACCEPT

# Masquerade outgoing traffic
iptables -t nat -I POSTROUTING -o upfgtp -j MASQUERADE
iptables -t nat -I POSTROUTING -o eth1 -j MASQUERADE

# Allow return traffic
iptables -I INPUT -m state --state RELATED,ESTABLISHED -j ACCEPT

# Disable reverse path filtering
for i in /proc/sys/net/ipv4/conf/*/rp_filter ; do
    echo 0 > $i 
done