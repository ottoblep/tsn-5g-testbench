#!/bin/bash
#
# Configure iptables in UE
#
# Forward everything
iptables-nft -I FORWARD -j ACCEPT

# Masquerade outgoing traffic
iptables-nft -t nat -I POSTROUTING -o oaitun_ue1 -j MASQUERADE
iptables-nft -t nat -I POSTROUTING -o eth1 -j MASQUERADE

# Allow return traffic
iptables-nft -I INPUT -m state --state RELATED,ESTABLISHED -j ACCEPT

# Disable reverse path filtering
for i in /proc/sys/net/ipv4/conf/*/rp_filter ; do
    echo 0 > $i 
done